# -*- coding: utf-8 -*-

from fabric.api import *


def install():
    """install"""
    local('go install github.com/itang/godict/cmd/godict')


def fmt():
    """format go code"""
    local('go fmt *.go')
    local('go fmt ./cmd/...')
