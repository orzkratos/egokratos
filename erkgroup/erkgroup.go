package erkgroup

import (
	"context"

	"github.com/orzkratos/errkratos"
	"golang.org/x/sync/errgroup"
)

type Group struct {
	ego *errgroup.Group
	ctx context.Context
}

func NewGroup(ctx context.Context) *Group {
	ego, ctx := errgroup.WithContext(ctx)
	return &Group{
		ego: ego,
		ctx: ctx,
	}
}

func (G *Group) Wait() *errkratos.Erk {
	if err := G.ego.Wait(); err != nil {
		return errkratos.FromError(err)
	}
	return nil
}

func (G *Group) Go(run func(ctx context.Context) *errkratos.Erk) {
	G.ego.Go(func() error {
		if erk := run(G.ctx); erk != nil {
			return erk
		}
		return nil
	})
}

func (G *Group) TryGo(run func(ctx context.Context) *errkratos.Erk) bool {
	return G.ego.TryGo(func() error {
		if erk := run(G.ctx); erk != nil {
			return erk
		}
		return nil
	})
}

func (G *Group) SetLimit(n int) {
	G.ego.SetLimit(n)
}
