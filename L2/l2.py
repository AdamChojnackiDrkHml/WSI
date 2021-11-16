import tensorflow as tf
import numpy as np

mnist = tf.keras.datasets.mnist

(training_data, training_labels), (test_data, test_labels) = mnist.load_data()
training_data, test_data = training_data / 255, test_data / 255

traing_data_falttened = training_data.reshape(len(training_data), 28*28)
test_data_flatened = test_data.reshape(len(test_data), 28*28)

model = tf.keras.Sequential([
    tf.keras.layers.Flatten(input_shape=(784,1)),
    tf.keras.layers.Dense(128, activation=tf.nn.relu),
    tf.keras.layers.Dense(10, activation=tf.nn.softmax)
])

model.compile(optimizer= tf.keras.optimizers.Adam(learning_rate=0.01),
              loss='sparse_categorical_crossentropy',
              metrics=['accuracy'])

model.fit(traing_data_falttened, training_labels, epochs=10)

model.evaluate(test_data_flatened, test_labels)


model.save('my_model')
