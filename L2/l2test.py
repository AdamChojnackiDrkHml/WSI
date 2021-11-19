from mlxtend.data import loadlocal_mnist
import tensorflow as tf
import numpy as np



model = tf.keras.models.load_model('my_model')
mnist = tf.keras.datasets.mnist
images_filename = 'my-images.idx3-ubyte'
label_filename = 'my-labels.idx1-ubyte'

train_images, train_labels = loadlocal_mnist(
            images_path=images_filename, 
            labels_path=label_filename)

model.evaluate(train_images, train_labels)
