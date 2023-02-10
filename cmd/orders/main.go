package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os/signal"

	. "applicationDesignTest/pkg/emails/usecase"
	. "applicationDesignTest/pkg/orders/endpoint/http"
	. "applicationDesignTest/pkg/orders/repository/inmemory"
	. "applicationDesignTest/pkg/orders/usecase"
	. "applicationDesignTest/pkg/rooms/usecase"

	"applicationDesignTest/utils/errgroup"
	"applicationDesignTest/utils/log/stdlog"
)

var AvailableRooms = map[string]struct{}{
	"econom":   {},
	"standart": {},
	"lux":      {},
}

const appName = "orderssvc"

func main() {
	ctx := context.Background()
	ctx, cancel := signal.NotifyContext(ctx)
	defer cancel()

	logger := stdlog.NewLogger(appName)
	logger.Infof("start app")
	defer logger.Infof("stop app")

	ordersRepository := NewOrdersRepository()

	emailsUsecase := NewEmailsUsecase()
	roomsUsecase := NewRoomsUsecase()
	ordersUsecase := NewOrdersUsecase(emailsUsecase, roomsUsecase, ordersRepository)

	ordersEntrypoint := NewOrdersEndpoint(ordersUsecase)

	mux := http.NewServeMux()
	ordersEntrypoint.Registrate(mux)

	httpServer := http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	wg, ctx := errgroup.WithContext(ctx)
	wg.Go(func() (err error) {
		err = httpServer.ListenAndServe()
		if errors.Is(err, http.ErrServerClosed) {
			err = nil
		} else if err != nil {
			err = fmt.Errorf("http server listen and serve: %w", err)
		}
		return
	})

	<-ctx.Done()
	logger.Infof("start application shutdown")
	httpServer.Shutdown(context.TODO())
	if err := wg.Wait(); err != nil {
		logger.Errorf("app running error: %s", err)
	}
}
