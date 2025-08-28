package tests

import (
    "testing"
    "supplx-gateway-marketplace/internal/lb"
)

func TestLeastConnections(t *testing.T){
    l := lb.NewLeastConnections([]string{"a","b","c"})
    a := l.Acquire()
    b := l.Acquire()
    c := l.Acquire()
    if a=="" || b=="" || c=="" { t.Fatal("expected endpoints") }
    if a==b || b==c { t.Fatal("expected distribution across endpoints") }
    l.Release(a); l.Release(b); l.Release(c)
}


