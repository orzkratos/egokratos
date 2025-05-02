package erkbizgroup

import (
	"context"

	"github.com/go-kratos/kratos/v2/errors"
	"github.com/orzkratos/errkratos/errbizkratos"
	"golang.org/x/sync/errgroup"
)

type Group struct {
	ego *errgroup.Group
}

func WithContext(ctx context.Context) (*Group, context.Context) {
	ego, ctx := errgroup.WithContext(ctx)
	return &Group{ego: ego}, ctx
}

func (G *Group) Wait() *errbizkratos.Ebz {
	if err := G.ego.Wait(); err != nil {
		return errbizkratos.NewEbz(errors.FromError(err))
	}
	return nil
}

func (G *Group) Go(run func() *errbizkratos.Ebz) {
	G.ego.Go(func() error {
		if ebz := run(); ebz != nil {
			return ebz.Erk
		}
		return nil
	})
}

func (G *Group) TryGo(run func() *errbizkratos.Ebz) bool {
	return G.ego.TryGo(func() error {
		if ebz := run(); ebz != nil {
			return ebz.Erk
		}
		return nil
	})
}

func (G *Group) SetLimit(n int) {
	G.ego.SetLimit(n)
}
