package safeconnmap

import (
    "net"
    "sync"
)

type SafeConnMap struct {
    mx sync.Mutex
    m map[string]net.Conn
}

func NewSafeConnMap() *SafeConnMap {
    return &SafeConnMap{ m: make(map[string]net.Conn) }
}

func (sm *SafeConnMap) Delete(key string) {
    sm.mx.Lock()
    defer sm.mx.Unlock()
    delete(sm.m, key)
}

func (sm *SafeConnMap) Add(key string, value net.Conn) {
    sm.mx.Lock()
    defer sm.mx.Unlock()
    sm.m[key] = value
}

func (sm *SafeConnMap) Read(key string) (net.Conn, bool) {
    val, ok := sm.m[key]
    return val, ok
}

func (sm *SafeConnMap) Exists(key string) (bool) {
    _, exists := sm.m[key]
    return exists
}

func (sm *SafeConnMap) Raw() (map[string]net.Conn) {
    return sm.m
}

