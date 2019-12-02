# SDF Path Tracer

A signed distance function path tracer, adapted from the ubiquitous *Ray Tracing in One Weekend* book, extended to include:

* various [2D](https://www.iquilezles.org/www/articles/distfunctions2d/distfunctions2d.htm) and [3D](http://iquilezles.org/www/articles/distfunctions/distfunctions.htm) SDFs
* SDF bounding spheres to allow fast(er) ray intersection and elimination
* multi-node cluster rendering via RPC

It seems pretty quick, at least in the ballpark of other similar efforts. The clustering feature seems less common; heaps of fun to spin up a few AWS bustable instances as render slaves and hammer lots of cores!

A shout-out to the author of [pt](https://github.com/fogleman/pt) is deserved. Very useful to have a Go-based reference implementation of a path tracer for debugging.

See [spt_test.go](spt_test.go) for this example:

![shapes.png](https://raw.githubusercontent.com/wiki/seanpringle/spt/shapes.png)