package tuple

import (
	"testing"

	"github.com/stretchr/testify/require"

	core "github.com/zapravila/spicedb/pkg/proto/core/v1"
)

func TestONRSet(t *testing.T) {
	set := NewONRSet()
	require.True(t, set.IsEmpty())
	require.Equal(t, uint32(0), set.Length())

	require.True(t, set.Add(ParseONR("resource:1#viewer")))
	require.False(t, set.IsEmpty())
	require.Equal(t, uint32(1), set.Length())

	require.True(t, set.Add(ParseONR("resource:2#viewer")))
	require.True(t, set.Add(ParseONR("resource:3#viewer")))
	require.Equal(t, uint32(3), set.Length())

	require.False(t, set.Add(ParseONR("resource:1#viewer")))
	require.True(t, set.Add(ParseONR("resource:1#editor")))

	require.True(t, set.Has(ParseONR("resource:1#viewer")))
	require.True(t, set.Has(ParseONR("resource:1#editor")))
	require.False(t, set.Has(ParseONR("resource:1#owner")))
	require.False(t, set.Has(ParseONR("resource:1#admin")))
	require.False(t, set.Has(ParseONR("resource:1#reader")))

	require.True(t, set.Has(ParseONR("resource:2#viewer")))
}

func TestONRSetUpdate(t *testing.T) {
	set := NewONRSet()
	set.Update([]*core.ObjectAndRelation{
		ParseONR("resource:1#viewer"),
		ParseONR("resource:2#viewer"),
		ParseONR("resource:3#viewer"),
	})
	require.Equal(t, uint32(3), set.Length())

	set.Update([]*core.ObjectAndRelation{
		ParseONR("resource:1#viewer"),
		ParseONR("resource:1#editor"),
		ParseONR("resource:1#owner"),
		ParseONR("resource:1#admin"),
		ParseONR("resource:1#reader"),
	})
	require.Equal(t, uint32(7), set.Length())
}

func TestONRSetIntersect(t *testing.T) {
	set1 := NewONRSet()
	set1.Update([]*core.ObjectAndRelation{
		ParseONR("resource:1#viewer"),
		ParseONR("resource:2#viewer"),
		ParseONR("resource:3#viewer"),
	})

	set2 := NewONRSet()
	set2.Update([]*core.ObjectAndRelation{
		ParseONR("resource:1#viewer"),
		ParseONR("resource:1#editor"),
		ParseONR("resource:1#owner"),
		ParseONR("resource:1#admin"),
		ParseONR("resource:2#viewer"),
		ParseONR("resource:1#reader"),
	})

	require.Equal(t, uint32(2), set1.Intersect(set2).Length())
	require.Equal(t, uint32(2), set2.Intersect(set1).Length())
}

func TestONRSetSubtract(t *testing.T) {
	set1 := NewONRSet()
	set1.Update([]*core.ObjectAndRelation{
		ParseONR("resource:1#viewer"),
		ParseONR("resource:2#viewer"),
		ParseONR("resource:3#viewer"),
	})

	set2 := NewONRSet()
	set2.Update([]*core.ObjectAndRelation{
		ParseONR("resource:1#viewer"),
		ParseONR("resource:1#editor"),
		ParseONR("resource:1#owner"),
		ParseONR("resource:1#admin"),
		ParseONR("resource:2#viewer"),
		ParseONR("resource:1#reader"),
	})

	require.Equal(t, uint32(1), set1.Subtract(set2).Length())
	require.Equal(t, uint32(4), set2.Subtract(set1).Length())
}

func TestONRSetUnion(t *testing.T) {
	set1 := NewONRSet()
	set1.Update([]*core.ObjectAndRelation{
		ParseONR("resource:1#viewer"),
		ParseONR("resource:2#viewer"),
		ParseONR("resource:3#viewer"),
	})

	set2 := NewONRSet()
	set2.Update([]*core.ObjectAndRelation{
		ParseONR("resource:1#viewer"),
		ParseONR("resource:1#editor"),
		ParseONR("resource:1#owner"),
		ParseONR("resource:1#admin"),
		ParseONR("resource:2#viewer"),
		ParseONR("resource:1#reader"),
	})

	require.Equal(t, uint32(7), set1.Union(set2).Length())
	require.Equal(t, uint32(7), set2.Union(set1).Length())
}

func TestONRSetWith(t *testing.T) {
	set1 := NewONRSet()
	set1.Update([]*core.ObjectAndRelation{
		ParseONR("resource:1#viewer"),
		ParseONR("resource:2#viewer"),
		ParseONR("resource:3#viewer"),
	})

	added := set1.With(ParseONR("resource:1#editor"))
	require.Equal(t, uint32(3), set1.Length())
	require.Equal(t, uint32(4), added.Length())
}

func TestONRSetAsSlice(t *testing.T) {
	set := NewONRSet()
	set.Update([]*core.ObjectAndRelation{
		ParseONR("resource:1#viewer"),
		ParseONR("resource:2#viewer"),
		ParseONR("resource:3#viewer"),
	})

	require.Equal(t, 3, len(set.AsSlice()))
}
