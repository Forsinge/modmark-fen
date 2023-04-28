# FEN-package for ModMark
This package is primarily meant as a proof of concept for package development in Go/TinyGo. It works by converting a FEN-string into an SVG.

## Usage
To use the package, simply paste a FEN-string as the contents to a [fen]-module. Example usage:
```
[fen width=0.6](r3k2r/p1ppqpb1/bn2pnp1/3PN3/1p2P3/2N2Q1p/PPPBBPPP/R3K2R w KQkq -)
```

## Building to WASM/WASI
To build locally with WASI as target, TinyGo can be used as follows: 
```
tinygo build -o fen.wasm -target wasi ./
```