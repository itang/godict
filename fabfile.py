# -*- coding: utf-8 -*-

from fabric.api import local 


def run():
    """run"""
    local('go run main.go')


def dev():
    """dev"""
    # https://github.com/tockins/realize
    local('realize run')


def repl():
    """repl"""
    local('gore')


def update():
    """dep ensure -update"""
    status()
    local('dep ensure -update')
    status()


def status():
    """dep status"""
    local('dep status')


def fmt():
    """go fmt ./..."""
    pkgs = ['github.com/itang/godict/cmd/godict', 'github.com/itang/godict']
    for pkg in pkgs:
        local('go fmt {}'.format(pkg))


def install():
    """install"""
    local('go install github.com/itang/godict/cmd/godict')
