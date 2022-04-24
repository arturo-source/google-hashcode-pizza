# google-hashcode-pizza
Solution for hashcode pizza problem 2017 and later

In this branch I try to solve the problem (still unsuccessful) by making a copy of the state of the pizza in each iteration, so I don't have to implement the logic of iterating over the slices.

The main problem doing by this way is that `input/medium.in` file gets overflow RAM, because of copying the pizza is too much. And `input/big.in` file makes the program work very slow, problably for the same reason.
