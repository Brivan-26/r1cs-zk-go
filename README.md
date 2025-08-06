# WIP

The following construction is not yet complete and is work-in-progress. Eventually, we will be there where we have a fully-working ZK construction.

> Example and reasoning are inspired from my journey reading [ZK-Book](https://rareskills.io/zk-book)

## Step1: (Not so) ZK for R1CS

> Implementation of step 1 is at commit [4aa7ad5](https://github.com/Brivan-26/r1cs-zk-go/tree/4aa7ad5951754495aa07378239b92b5352c61f25)

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

## Step 2: Making the proof system succinct, building QAP
In the previous step, we established a foundation for our proving system. However, a significant issue remains: the construction is not yet succinct. Specifically, if we have $n$ constraints, the prover must send $3n$ points, and the verifier must perform $n$ pairings and checks.

To address this, let's reconsider what the verifier is checking: $La * Ra = Oa$. Here, the verifier is comparing two vectors, which requires $O(n)$ complexity. Interestingly, comparing two polynomials is much more efficient, with a complexity of $O(1)$, thanks to the [Schwartz-Zippel Lemma](https://en.wikipedia.org/wiki/Schwartz%E2%80%93Zippel_lemma). To leverage this, we must convert our system from a vector-based to a polynomial-based representation. This is known as a Quadratic Arithmetic Program (QAP). We cannot convert the entire matrix $L$ to a single polynomial at once, so instead, we break it into multiple vectors—one for each column—and interpolate each column. This is because:
```math 
\begin{bmatrix} 0 & 0 & 1 & 0 \\ -5 & 1 & 0 & -5 \end{bmatrix} \begin{bmatrix} 1 \\ y \\ v \\ x \end{bmatrix} = \begin{bmatrix} 0 \\ -5 \end{bmatrix} * 1 + \begin{bmatrix} 0 \\ 1 \end{bmatrix} * y + \begin{bmatrix} 1 \\ 0 \end{bmatrix} * v + \begin{bmatrix} 0 \\ -5 \end{bmatrix} * x
```
Therefore, for the $L$ matrix (and similarly for $R$ and $O$), we interpolate each column as $u_i(x)$ and multiply it by its corresponding element from the witness vector. Formally, for the example above:
```math
La =  u_1(x) + y u_2(x) + v u_3(x) + x u_4(x) = \sum_{i=1}^{4} a_i u_i(x) = u(x)
```
```math
Ra =  v_1(x) + y v_2(x) + v v_3(x) + x v_4(x) = \sum_{i=1}^{4} a_i v_i(x) = v(x)
```
```math
Oa =  w_1(x) + y w_2(x) + v w_3(x) + x w_4(x) = \sum_{i=1}^{4} a_i w_i(x) = w(x)
```
This works because **The group of vectors under addition in a finite field is homomorphic to the group of polynomials under addition in a finite field.**

Naturally, the degrees of $u(x)$, $v(x)$, and $w(x)$ will be at most $n-1$, where $n$ is the number of constraints (rows).

Our proof system therefore becomes:
```math
u(x) * v(x) = w(x)
```
However, there is an issue: the polynomial $u(x) * v(x)$ will have degree $2n-2$ and will not generally equal the interpolated $w(x)$ (which is of degree $n-1$), since the homomorphism we established applies to addition and scalar multiplication, not the Hadamard product. To address this, we balance the equation as follows:
```math
u(x) * v(x) = w(x) + b(x)
```
We must constrain the polynomial $b(x)$ to be interpolated from the zero vector, so we do not invalidate our proof. This is analogous to adding the zero vector to the underlying vectors being interpolated: $v_0 * v_1 = v + 0$. Thus, we can force $b(x)$ to be factored by another polynomial that has roots at the points $0, 1, ..., n$—the interpolation points—denoted as $t(x)$, so $b(x) = t(x) * h(x)$, where $t(x) = x(x-1)...(x-n)$.

Our proof system thus becomes:
```math
u(x) * v(x) = w(x) + t(x)*h(x)
```
```math
\sum_{i=1}^{4} a_i u_i(x) * \sum_{i=1}^{4} a_i v_i(x) = \sum_{i=1}^{4} a_i w_i(x) + t(x)h(x)
```
$h(x)$ can be easily computed from the other polynomials in the equation above. Now, the verifier can choose a random number $r$, and the prover responds with:
```math
A = u_x(r), B = v_x(r), C = w_x(r) + t_x(r)h_x(r)
```
The verifier then simply checks that:
```math
AB = C
```
Thanks to the [Schwartz-Zippel Lemma](https://en.wikipedia.org/wiki/Schwartz%E2%80%93Zippel_lemma), this ensures with very high probability that the two polynomials are equal and thus, the underlying vectors interpolated are equal!.
However, this requires the verifier to trust that the prover is evaluating the polynomials correctly. If the prover knows the point $r$, they could fabricate points that satisfy the above condition. To resolve this, the prover must evaluate the polynomials at a secret value, $\tau$. This is the basis for the **Trusted Setup Ceremony**.

Evaluating polynomials simply involves multiplying coefficients by the evaluation point:
```math
f(x) = \langle [c_n, ..., c_3, c_2, c_1, c_0],\ [x^n, ..., x^3, x^2, x, 1] \rangle
```

We can generate a secret value $\tau$ ahead of time, encrypt it by multiplying with $G_1$, and then send a vector consisting of $\tau$ raised to successive degrees. This is known as the **Structured Reference String (SRS)**: $G_1, \tau * G_1, \tau^2 * G_1, ..., \tau^n * G_1$, where $n$ is the degree of the polynomial.

Since we require multiplication, we also generate another SRS where $\tau$ is encrypted in $G_2$:
```math
SRS1 = [G_1, G_1\tau, G_1\tau^2, ..., G_1\tau^n]
```
```math
SRS2 = [G_2, G_2\tau, G_2\tau^2, ..., G_2\tau^n] \\ 
```
$u_x$ will be evaluated using SRS1, $v_x$ using SRS2, and $w_x$ using SRS1. However, we need to evaluate $t_xh_x$. The previous method does not work here, because $h_x$ is known only to the prover and is not available during the Trusted Setup ceremony. However, we know that $t(x)h(x)$ evaluated at a point $r$ is $t(r)h_x(r)$. Therefore, during the Trusted Setup, we can calculate $SRS3$, which consists of:
```math
SRS3 = [t(r)G_1, G_1t(r)\tau, G_1t(r)\tau^2, ..., G_1t(r)\tau^(n-2)]
```
This works because the polynomial $t(x)$ is publicly known.

Therefore, the final output of the Trusted Setup Ceremony is:
```math
SRS1 = [G_1, G_1\tau, G_1\tau^2, ..., G_1\tau^n] = [\Omega_{n-1},\, \Omega_{n-2},\, \ldots,\, \Omega_{1},\, G_{1}]
```
```math
SRS2 = [G_2, G_2\tau, G_2\tau^2, ..., G_2\tau^n] = [\Theta_{n-1},\, \Theta_{n-2},\, \ldots,\, \Theta_{1},\, G_{2}]
```
```math
SRS3 = [t(r)G_1, G_1t(r)\tau, G_1t(r)\tau^2, ..., G_1t(r)\tau^(n-2)] = [\Upsilon_{n-2},\, \Upsilon_{n-3},\, \ldots,\, \Upsilon_{1},\, \Upsilon_{0}]
```

Now, the prover can efficiently compute $A$, $B$, and $C$ using the $SRS$:
```math
\begin{aligned}
A &= \sum_{i=1}^{m} a_i u_i(\tau) = \langle [u_{n-1}, u_{n-2}, \ldots, u_{1}, u_{0}],\ [\Omega_{n-1}, \Omega_{n-2}, \ldots, \Omega_1, G_1] \rangle \\[1em]
B &= \sum_{i=1}^{m} a_i v_i(\tau) = \langle [v_{n-1}, v_{n-2}, \ldots, v_{1}, v_{0}],\ [\Theta_{n-1}, \Theta_{n-2}, \ldots, \Theta_1, G_2] \rangle \\[1em]
C &= \sum_{i=0}^{m} a_i w_i(\tau) + h(\tau) t(\tau) = \langle [w_{n-1}, w_{n-2}, \ldots, w_{1}, w_{0}],\ [\Omega_{n-1}, \Omega_{n-2}, \ldots, \Omega_1, G_1] \rangle \\ 
&\hspace{3cm} + \langle [h_{n-2}, h_{n-3}, \ldots, h_1, h_0],\ [\Upsilon_{n-2}, \Upsilon_{n-3}, \ldots, \Upsilon_1, \Upsilon_0] \rangle
\end{aligned}
```
The prover publishes $A$, $B$, and $C$, and the verifier simply checks the pairing:
```math
pairing(A, B) == pairing(C, G2)
```

With this, we have achieved succinctness! Regardless of the number of constraints, the prover sends only three points, and the verifier needs to verify only a single pairing!

The implementation of this step is available at the [HEAD commit](https://github.com/Brivan-26/r1cs-zk-go)

### What we need to solve in next steps
1. Given that the verifier checks only the pairing of three points, the prover could send arbitrary points $A$, $B$, $C$ that satisfy the pairing, and the verifier has no guarantee that these points were derived from the generated QAP.
2. Our proof system does not yet support public inputs, so the verifier has no means of injecting public inputs.