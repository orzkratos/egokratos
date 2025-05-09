package ebzgroup

import (
	"context"

	"github.com/go-kratos/kratos/v2/errors"
	"github.com/orzkratos/errkratos/ebzkratos"
	"golang.org/x/sync/errgroup"
)

type Group struct {
	ego *errgroup.Group
}

func WithContext(ctx context.Context) (*Group, context.Context) {
	ego, ctx := errgroup.WithContext(ctx)
	return &Group{ego: ego}, ctx
}

func (G *Group) Wait() *ebzkratos.Ebz {
	if err := G.ego.Wait(); err != nil {
		return ebzkratos.NewEbz(errors.FromError(err))
	}
	return nil
}

func (G *Group) Go(run func() *ebzkratos.Ebz) {
	G.ego.Go(func() error {
		if ebz := run(); ebz != nil {
			return ebz.Erk
		}
		return nil
	})
}

func (G *Group) TryGo(run func() *ebzkratos.Ebz) bool {
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
