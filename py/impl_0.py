# ----------------------------------------------------------------------
# Ref. implementation example:
# using BN254 curve with ATE pairing
# ----------------------------------------------------------------------

from bn254.curve import r as scalar_upper_bound
from bn254.ecp import generator as generator1 
from bn254.ecp2 import generator as generator2
from bn254.big import invmodp, rand
from bn254.pair import e

# --------------- Private key calculation ------------------------------

#recipient side
k, v = rand(scalar_upper_bound), rand(scalar_upper_bound)

K, V = k * generator2(),  v * generator1()


#sender side
r = rand(scalar_upper_bound)

R = r * generator1()

# --------------- Recipient's Public key `P` calculation ---------------


#recipient side

#sender side

P = e(K, r*V) 


P2 = k * v * generator2()

print(P.toBytes(), '\n\n\n', P2)

# --------------- end. --------------------------------------------------

