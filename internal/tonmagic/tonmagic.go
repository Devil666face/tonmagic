package tonmagic

import (
	"crypto/tls"
	"fmt"
	"net"
	"net/url"
	"strings"
	"tonmagic/internal/config"
	"tonmagic/pkg/lib"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/helmet"
	"github.com/gofiber/fiber/v3/middleware/limiter"
	"github.com/gofiber/fiber/v3/middleware/logger"
	"github.com/gofiber/fiber/v3/middleware/proxy"
	"github.com/gofiber/fiber/v3/middleware/recover"
	"github.com/pkg/errors"
	"github.com/valyala/fasthttp"

	ton "github.com/xssnick/tonutils-proxy/proxy"
)

const adnl = "adnl"

var listenConfig = fiber.ListenConfig{DisableStartupMessage: true}

type Tonmagic struct {
	tonlisten string
	config    *config.Config
	tlsconfig *tls.Config
	httpApp   *fiber.App
	httpsApp  *fiber.App
}

func New(_config *config.Config, _tlsconfig *tls.Config) (*Tonmagic, error) {
	_tonlisten, err := lib.RandomLocalConnect()
	if err != nil {
		return nil, errors.Wrap(err, "failed to init address for ton proxy")
	}
	t := Tonmagic{
		tonlisten: _tonlisten,
		tlsconfig: _tlsconfig,
		config:    _config,
	}
	t.httpApp = t.NewApp()
	t.httpsApp = t.NewApp()
	return &t, nil
}

func (t *Tonmagic) NewApp() *fiber.App {
	app := fiber.New()

	app.Use(logger.New())
	app.Use(recover.New())
	app.Use(helmet.New(helmet.Config{
		XSSProtection:  "1",
		ReferrerPolicy: "same-origin",
	}))
	app.Use(limiter.New(limiter.Config{
		Max: 100000,
	}))
	app.Use(func(c fiber.Ctx) error {
		proxy.WithClient(&fasthttp.Client{
			Dial: func(addr string) (net.Conn, error) {
				return fasthttp.Dial(t.tonlisten)
			},
		})
		c.Request().Header.Set("X-Forwarded-Proto", "http")
		host, err := SplitHost(c.Host())
		if err != nil {
			return err
		}
		fmt.Println(host)
		if err := proxy.Do(c, fmt.Sprintf("http://%s/%s", host, c.OriginalURL())); err != nil {
			return err
		}
		c.Response().Header.Del("Server")
		c.Response().Header.Del("Ton-Reverse-Proxy")
		return nil
	})

	return app
}

func (t *Tonmagic) StartTonProxy() error {
	if _, err := ton.StartProxy(t.tonlisten, false, nil, "CLI ", false); err != nil {
		return errors.Wrap(err, "failed to start proxy")
	}
	return nil
}

func (t *Tonmagic) ListenHttp() error {
	if err := t.httpApp.Listen(t.config.HttpConnect, listenConfig); err != nil {
		return errors.Wrap(err, "error listen http")
	}
	return nil
}

func (t *Tonmagic) ListenHttps() error {
	lis, err := net.Listen("tcp", t.config.HttpsConnect)
	if err != nil {
		return errors.Wrap(err, "error to listen")
	}

	if err := t.httpsApp.Listener(tls.NewListener(lis, t.tlsconfig), listenConfig); err != nil {
		return errors.Wrap(err, "error listen https")
	}
	return nil
}

func SplitHost(host string) (string, error) {
	u, err := url.Parse("//" + host)
	if err != nil {
		return "", errors.Wrap(err, "error parse url")
	}

	hostParts := strings.Split(u.Hostname(), ".")
	if len(hostParts) != 4 || hostParts[1] != adnl {
		return "", errors.New("failed to parse adnl url")
	}
	return strings.Join(hostParts[:len(hostParts)-2], "."), nil
}
