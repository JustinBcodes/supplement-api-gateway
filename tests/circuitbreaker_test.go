package tests

import (
    "testing"
    "time"
    "supplx-gateway-marketplace/pkg/breaker"
)

func TestBreakerTransitions(t *testing.T){
    b := breaker.New(5, 0.5, 50*time.Millisecond, 2, 10)
    for i:=0;i<5;i++{ b.OnResult(false) }
    if b.State() == breaker.Closed { t.Fatal("expected not closed") }
}


