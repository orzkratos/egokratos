package erkgroup

import (
	"context"

	"github.com/go-kratos/kratos/v2/errors"
	"github.com/orzkratos/errkratos"
	"golang.org/x/sync/errgroup"
)

type Group struct {
	ego *errgroup.Group
}

func WithContext(ctx context.Context) (*Group, context.Context) {
	ego, ctx := errgroup.WithContext(ctx)
	return &Group{ego: ego}, ctx
}

func (G *Group) Wait() *errkratos.Erk {
	if err := G.ego.Wait(); err != nil {
		return errors.FromError(err)
	}
	return nil
}

func (G *Group) Go(run func() *errkratos.Erk) {
	G.ego.Go(func() error {
		if erk := run(); erk != nil {
			return erk
		}
		return nil
	})
}

func (G *Group) TryGo(run func() *errkratos.Erk) bool {
	return G.ego.TryGo(func() error {
		if erk := run(); erk != nil {
			return erk
		}
		return nil
	})
}

func (G *Group) SetLimit(n int) {
	G.ego.SetLimit(n)
}
