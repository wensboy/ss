package server

import (
	"context"
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v5"
	"github.com/urfave/cli/v3"
	"github.com/wensboy/ss/config"
)

func mockMiddleware(name string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c *echo.Context) error {
			fmt.Printf("hit [%s] mock middleware\n", name)
			return next(c)
		}
	}
}

func MockMount(module string) RestServerHook {
	return func(s *RestServer) {
		fmt.Printf("mount %s\n", module)
	}
}

func MockLoadAll() {
	cfg := config.NewConfig().SetDir("../spec").SetExt("json")
	cfg.Load(
		"config",
		"server",
		"db",
	)
	fmt.Printf("%#v\n", cfg)
	config.SetGlobalConfig(cfg)
	err := godotenv.Load("../spec/.env")
	if err != nil {
		panic(err)
	}
}

var (
	mockMountConfig = func(s *RestServer) {
		fmt.Printf("mount config\n")
		s.server.Addr = ":18080"
	}
	mockMountDatabases = MockMount("databases")
	mockRegisterRouter = func(g *echo.Group, _ *ServerContext) {
		poblicGroup := g.Group("/public")
		poblicGroup.GET("/ping", func(c *echo.Context) error {
			return c.String(http.StatusOK, "pong")
		})
	}
)

func Test_GoStart(t *testing.T) {
	s := NewRestServer()
	s.MountModules(
		s.MountDatabases(nil, mockMountDatabases, nil),
		s.MountConfig(nil, nil, nil),
		s.MountMiddlewares("", mockMiddleware("auth")),
		s.MountRouters("", mockRegisterRouter),
	)
	ctx, cancel := context.WithCancel(context.Background())
	s.GoStart(ctx)
	time.Sleep(20 * 1_000 * time.Millisecond) // 20s 验证 router 注册
	cancel()
}

func Test_MountConfig(t *testing.T) {
	MockLoadAll()
	rootCmd := config.InitCommand("../spec/command.json")
	serveCmd := config.GetCommand("ss.serve")
	serveCmd.Action = func(ctx context.Context, c *cli.Command) error {
		s := NewRestServer()
		s.MountModules(
			s.MountDatabases(nil, mockMountDatabases, nil),
			s.MountConfig(nil, nil, nil),
			s.MountMiddlewares("", mockMiddleware("auth")),
			s.MountRouters("", mockRegisterRouter),
		)
		ctx, cancel := context.WithCancel(context.Background())
		s.GoStart(ctx)
		time.Sleep(20 * 1_000 * time.Millisecond) // 20s 验证 router 注册
		cancel()
		return nil
	}
	if err := rootCmd.Run(context.TODO(), []string{"ss", "serve"}); err != nil {
		t.Error(err)
	}
}
