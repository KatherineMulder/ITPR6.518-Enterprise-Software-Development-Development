package main

import (
    "testing"
    "os"
)

func TestMain(m *testing.M) {
    // Call flag.Parse() here if TestMain uses flags
    os.Exit(m.Run())
}
