package main

import (
	"math/rand"
	"testing"
)

func TestAxrot3d(t *testing.T) {
	for r := axrot3d(0); r < 24; r++ {
		x := coord3i{rand.Intn(7) - 3, rand.Intn(7) - 3, rand.Intn(7) - 3}
		x2 := r.rotate(x)
		x3 := r.inverse().rotate(x2)
		if x != x3 {
			t.Errorf("%s.inverse(%s.rotate(%s))=%s", r, r, x, x3)
			t.Errorf("reported inverse of %s is %s", r, r.inverse())
		}
	}
}

func TestAxform3d(t *testing.T) {
	for r := axrot3d(0); r < 24; r++ {
		tr := coord3i{rand.Intn(21) - 10, rand.Intn(21) - 10, rand.Intn(21) - 10}
		f := axform3d{r: r, t: tr}
		x := coord3i{rand.Intn(7) - 3, rand.Intn(7) - 3, rand.Intn(7) - 3}
		x2 := f.apply(x)
		x3 := f.inverse().apply(x2)
		if x != x3 {
			t.Errorf("%s.inverse(%s.apply(%s))=%s", f, f, x, x3)
		}
	}
}

func TestAxrot3dCompose(t *testing.T) {
	for r1 := axrot3d(0); r1 < 24; r1++ {
		for r2 := axrot3d(0); r2 < 24; r2++ {
			x := coord3i{rand.Intn(7) - 3, rand.Intn(7) - 3, rand.Intn(7) - 3}
			want := r2.rotate(r1.rotate(x))
			got := (r1.compose(r2)).rotate(x)
			if want != got {
				t.Errorf("%s.compose(%s).rotate(%s) = %s, want %s", r1, r2, x, got, want)
			}
		}
	}
}

func TestAxform3dCompose(t *testing.T) {
	for cases := 0; cases < 1000; cases++ {
		r1 := axrot3d(rand.Intn(24))
		tr1 := coord3i{rand.Intn(21) - 10, rand.Intn(21) - 10, rand.Intn(21) - 10}
		r2 := axrot3d(rand.Intn(24))
		tr2 := coord3i{rand.Intn(21) - 10, rand.Intn(21) - 10, rand.Intn(21) - 10}
		f1 := axform3d{r: r1, t: tr1}
		f2 := axform3d{r: r2, t: tr2}
		x := coord3i{rand.Intn(7) - 3, rand.Intn(7) - 3, rand.Intn(7) - 3}
		want := f2.apply(f1.apply(x))
		got := (f1.compose(f2)).apply(x)
		if want != got {
			t.Errorf("%s.compose(%s)(%s) = %s, want %s", f1, f2, x, got, want)
			t.Errorf("we got %s.compose(%s) = %s", f1, f2, f1.compose(f2))
		}
	}
}
