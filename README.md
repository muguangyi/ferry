# Seek

[![Build Status](https://travis-ci.com/muguangyi/seek.svg?branch=master)](https://travis-ci.com/muguangyi/seek) [![GoDoc](https://godoc.org/github.com/muguangyi/seek?status.svg)](https://godoc.org/github.com/muguangyi/seek) [![codecov](https://codecov.io/gh/muguangyi/seek/branch/master/graph/badge.svg)](https://codecov.io/gh/muguangyi/seek) [![Go Report Card](https://goreportcard.com/badge/github.com/muguangyi/seek)](https://goreportcard.com/report/github.com/muguangyi/seek)

**Seek** is a lightweight server develop framework for `golang`. **Seek** setup connections between containers based on `feature dependency`. This solution provides much flexibility for user to setup extendable servers quickly.

## What Seek DO

**Seek** is a server framework, and define the dev pattern to standalize server startup, connection and communicate flow. So user DO NOT need to write code for low level logic, like network connection, communicate protocol, etc, but only focus on the design and implementation for internal modules. It could make the modules more cohesive, decomposed, and general to improve the reusability.

## What Seek DO NOT do

There is no `gateway`, `lobby`, or `login` server implementation in **Seek**, even no `log` module. Those featured modules will not be provided by **Seek**, but need user to implement based on **Seek** framework.

## Framework Diagram

    +--------------------+            +=======+  register  +---------------+
    | dock               |  register  |       |<<<<<<>>>>>>| dock          |
    |                    |<<<<<<>>>>>>|       |   query    +---------------+
    |                    |   query    |  hub  |
    |                    |            |       |  register  +---------------+
    | +----------------+ |            |       |<<<<<<>>>>>>| dock          |
    | | feature 1      | |            +=======+   query    |               |
    | | feature 2      | |                                 | +-----------+ |
    | | book feature N | |<------------------------------->| | feature N | |
    | +----------------+ |        directly connected       | +-----------+ |
    +--------------------+                                 +---------------+

## Tech Notes

* Feature container (dock) is an independent server node, and could contain many features.
* Every feature runs within an independent routine.
* Communication between features base on channel RPC (only support sync mode so far)
* Features in different docks could communicate through the same way (RPC based on `feature dependency`)

## Limitation

Can't pass `func` or `interface` as parameter to feature's export methods. That's because the communication between features should be **data**, but not logic since **Seek** can't make a shadow for func or interface and do the data transition between different features.

So the **RULE** for feature design is: **DO NOT** define `func` or `interface` as method parameter!!!

## Quick Start

## Document