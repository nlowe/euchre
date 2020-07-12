package hosting

import (
	"context"

	"github.com/gammazero/nexus/v3/client"
	"github.com/gammazero/nexus/v3/router"
	"github.com/gammazero/nexus/v3/wamp"
	"github.com/sirupsen/logrus"
)

type TableServer struct {
	client *client.Client
}

func NewTableServer(nxr router.Router, table string) (*TableServer, error) {
	subscriber, err := client.ConnectLocal(nxr, client.Config{
		Realm:  table,
		Logger: logrus.WithFields(logrus.Fields{"prefix": "wamp::client", "table": table}),
	})
	if err != nil {
		return nil, err
	}

	err = subscriber.Register("say", func(ctx context.Context, invocation *wamp.Invocation) client.InvokeResult {
		if len(invocation.Arguments) < 2 {
			return client.InvokeResult{
				Args:   wamp.List{"Invalid Number of Arguments"},
				Kwargs: wamp.Dict{"got": len(invocation.Arguments), "want at least": 2},
				Err:    "error.syntax",
			}
		}

		if err := subscriber.Publish("state", wamp.Dict{}, invocation.Arguments, wamp.Dict{}); err != nil {
			return client.InvokeResult{
				Args: wamp.List{"Failed to publish state to other players"},
				Err:  "error.server",
			}
		}

		return client.InvokeResult{Args: wamp.List{"ok"}}
	}, wamp.Dict{"disclose_caller": true})
	if err != nil {
		return nil, err
	}

	return &TableServer{client: subscriber}, nil
}

func (t *TableServer) Close() error {
	return t.client.Close()
}
