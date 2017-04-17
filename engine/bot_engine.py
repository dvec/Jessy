#sudo pip install https://github.com/pybrain/pybrain/archive/0.3.3.zip
from pybrain.structure import LinearLayer, SigmoidLayer, FullConnection, RecurrentNetwork
from pybrain.supervised.trainers import BackpropTrainer
from pybrain.datasets import SupervisedDataSet

n = RecurrentNetwork()

n.addInputModule(LinearLayer(2, name='in'))
n.addModule(SigmoidLayer(3, name='hidden'))
n.addOutputModule(LinearLayer(1, name='out'))

n.addConnection(FullConnection(n['in'], n['hidden'], name='c1'))
n.addConnection(FullConnection(n['hidden'], n['out'], name='c2'))
n.addRecurrentConnection(FullConnection(n['hidden'], n['hidden'], name='c3'))

n.sortModules()

ds = SupervisedDataSet(2, 1)
ds.addSample((0, 0), (0,))
ds.addSample((0, 1), (1,))
ds.addSample((1, 0), (1,))
ds.addSample((1, 1), (0,))

trainer = BackpropTrainer(n, ds)
trainer.trainUntilConvergence()

print(n.activate([0, 0]))
