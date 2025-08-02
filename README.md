# WIP

The following construction is not yet complete and is work-in-progress. Eventually, we will be there where we have a fully-working ZK construction.

> Example and reasoning are inspired from my journey reading [ZK-Book](https://rareskills.io/zk-book)

## Step1: (Not so) ZK for R1CS

It is so brilliant how far you can go in building zero-knowledge with just ECC points and pairings. This repository is an ongoing implementation of ZK construction and it supports any problem statement that can be written as a set of constraints (circuit). You only need to adapt the circuit matrices ($L$, $R$, $O$ in `main.go`). 

The code in [main.go](./main.go) is a (not so) zk code that proves the following statement:
> I know a number `x` which is a solution for $x^3 + 5x + 5 = 155$

We start first by building our arithmetic circuit in an R1CS format:
```math
v = x^2
```
```math
y = x*v + 5x + 5
```
To make our life easier, we convert this to the following R1CS to have it in the form $O*a = La * Ra$:
```math
v = x^2
```
```math
y - 5x - 5 = x*v
```
Our witness vector is $a = [1, y, v, x]$

So, the matrix representation for the R1CS is the following:
```math 
\begin{bmatrix} 0 & 0 & 1 & 0 \\ -5 & 1 & 0 & -5 \end{bmatrix} \begin{bmatrix} 1 \\ y \\ v \\ x \end{bmatrix} = \begin{bmatrix} 0 & 0 & 0 & 1 \\ 0 & 0 & 0 & 1 \end{bmatrix} \begin{bmatrix} 1 \\ y \\ v \\ x \end{bmatrix} * \begin{bmatrix} 0 & 0 & 0 & 1 \\ 0 & 1 & 0 & 0 \end{bmatrix} \begin{bmatrix} 1 \\ y \\ v \\ x \end{bmatrix} 
```

The naive approach is to send the vector $a$ to the verifier, and the latter verifies that the equation holds. But, of course, this is not zero-knowledge.

So, instead, we encrypt the vector $a$. However, we still need to be able to perform operations on the encrypted values, and since the elliptic curve point group is homomorphic to the field for the addition operation, it is well suited for this. Given that we are performing the Hadamard product above, we will need to multiply elliptic curve points at some point, which requires the usage of pairings. That's why we need to encrypt the left side $La$ using $G_1$ and the right side $Ra$ using $G_2$. As a result, our encrypted vector $a$ for $G_1$ and $G_2$ is:
```math
\begin{bmatrix} G1 \\ yG1 \\ vG1 \\ xG1 \end{bmatrix} - 
\begin{bmatrix} G2 \\ yG2 \\ vG2 \\ xG2 \end{bmatrix}

```
Given the hardness of the discrete log problem, a verifier cannot, for example, determine $x$ from the point $xG_1$.

So, the matrix representation for the R1CS becomes:
```math 
\begin{bmatrix} 0 & 0 & 1 & 0 \\ -5 & 1 & 0 & -5 \end{bmatrix} \begin{bmatrix} 1 \\ y \\ v \\ x \end{bmatrix} = \begin{bmatrix} 0 & 0 & 0 & 1 \\ 0 & 0 & 0 & 1 \end{bmatrix} \begin{bmatrix} G1 \\ yG1 \\ vG1 \\ xG1 \end{bmatrix} * \begin{bmatrix} 0 & 0 & 0 & 1 \\ 0 & 1 & 0 & 0 \end{bmatrix} \begin{bmatrix} G2 \\ yG2 \\ vG2 \\ xG2 \end{bmatrix} 
```
After performing the Hadamard product, we get:
```math
\begin{bmatrix} v \\ -5 + y -5x \end{bmatrix} = \begin{bmatrix} xG1 \\ xG1 \end{bmatrix} * \begin{bmatrix} xG2 \\ yG2 \end{bmatrix}
```
Notice that we haven't encrypted the witness $a$ on the output side $Oa$. The naive approach would be to encrypt the values in the group $G_{12}$, since the pairing of a point from $G_1$ ($La$) and another point from $G_2$ ($Ra$) will result in a point in the target group $G_{12}$. But exchanging points in the group $G_{12}$ is extremely impractical since the latter is encoded in 12 dimensions. However, we can rely on the following:
```math
e(laG1, raG2) = e(oaG1, G2)
```
This comes from the bilinearity property, more specifically:
```math
e(laG1, raG2) = e(oaG1, G2) \\
=> e(G1, G2)^{l a^{2} r} = e(G1, G2)^{oa} \\
=> la^{a}r = oa \\
=> lar = o
```
So, the prover needs to encrypt $Oa$ in $G_1$. by doing so, we get the following R1CS:
```math
\begin{bmatrix} vG1 \\ -5G1 + yG1 -5xG1 \end{bmatrix} = \begin{bmatrix} xG1 \\ xG1 \end{bmatrix} * \begin{bmatrix} xG2 \\ yG2 \end{bmatrix}
```
The correct solution is $x = 5$. So, substituting $y$, $x$, and $v$ we get:
```math
\begin{bmatrix} 25G1 \\ -5G1 + 155G1 -25G1 \end{bmatrix} = \begin{bmatrix} 5G1 \\ 5G1 \end{bmatrix} * \begin{bmatrix} 5G2 \\ 155G2 \end{bmatrix}
```

Now, the prover sends the above 6 points to the verifier and the latter calculates and checks the pairing result. More specifically, the verifier checks that:
1. $e(5G1, 5G2) == e(25G1, G2)$
2. $e(5G1, 155G2) == e(125G1, G2)$

If both the above verifications pass, then the prover has sent a valid proof and he knows the number $x$.

### What we need to solve in next steps
1. The construction is not zero-knowledge secure. A verifier, by doing guesses, can infer the witness $a$ by simply multiplying his guess with $G_1$ and $G_2$ and comparing the result with the points sent by the prover.
2. Definitely, the construction is not succinct. The verifier needs to do two pairings in our example to check the proof. We want to have a $O(1)$ verification complexity
