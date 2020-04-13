import numpy as np
import rrcf

# A (robust) random cut tree can be instantiated from a point set (n x d)
X = np.random.randn(10, 2)
tree = rrcf.RCTree(X)
print(tree)


tree.insert_point([1, 2], 'random')

print(tree)
print(tree.codisp('random'))
