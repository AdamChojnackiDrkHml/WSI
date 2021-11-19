# Load, normalize and save every image in input_path folder
# into a file with similiar format to MNIST
# First character in image name will be saved as label for that file
# eg. 1asdasd.png will have label equal to 1

import os
from PIL import Image
import numpy as np
from scipy import ndimage

input_path = './My Set/'
image_filename = 'my-images2.idx3-ubyte'
labels_filename = 'my-labels2.idx1-ubyte'

def getBoundingBox(image):
    top_x = image.size
    top_y = image.size
    bot_x = 0
    bot_y = 0

    current_row = 0
    current_column = 0
    for row in image:
        current_column = 0
        for pix in row:
            if pix != 0:
                if current_row < top_y:
                    top_y = current_row
                if current_column < top_x:
                    top_x = current_column
                if(current_column > bot_x):
                    bot_x = current_column
                if(current_row > bot_y):
                    bot_y = current_row
                

            current_column += 1
        current_row += 1

    return top_x, top_y, bot_x, bot_y

def parseFile(filepath):
    img_data = Image.open(input_path + filename)
    # Greyscale
    gs_image = img_data.convert(mode='L')
    # Cut image to number bounding box
    pixel_array = np.asarray(gs_image)
    bounding_box = getBoundingBox(pixel_array)
    pixel_array = pixel_array[bounding_box[1]:bounding_box[3],bounding_box[0]:bounding_box[2]]
    # Scale image to 20x20
    gs_image = Image.fromarray(pixel_array, 'L')
    gs_image.thumbnail([20,20],Image.ANTIALIAS)
    # Calculate center of mass for pixels
    pixel_array = np.asarray(gs_image)
    cy, cx = ndimage.measurements.center_of_mass(pixel_array)
    center = (int(cx), int(cy))
    # Insert 20x20 image into 28x28 (using center of mass of pixels as middle point)
    blank = np.zeros((28,28),dtype=np.byte)
    x_offset = int((13 - center[0]))
    y_offset = int((13 - center[1]))

    if x_offset + pixel_array.shape[1] > 27:
        x_offset = 27 - pixel_array.shape[1]
    if y_offset + pixel_array.shape[0] > 27:
        y_offset = 27 - pixel_array.shape[0]

    index_x = 0
    index_y = 0
    while index_y < pixel_array.shape[0]:
        while index_x < pixel_array.shape[1]:
            blank[index_y+y_offset][index_x+x_offset] = pixel_array[index_y][index_x]
            index_x += 1
        index_y += 1
        index_x = 0
    
    return blank.flatten()

files = os.listdir(input_path)

# Create binary output files
images_file = open(image_filename,'wb')
labels_file = open(labels_filename,'wb')

# Add headers
images_file.write(b"\x00\x00\x08\x01")
labels_file.write(b"\x00\x00\x08\x03")
images_file.write((len(files)).to_bytes(4, byteorder='big'))
labels_file.write((len(files)).to_bytes(4, byteorder='big'))
images_file.write((28).to_bytes(4, byteorder='big'))
images_file.write((28).to_bytes(4, byteorder='big'))
print("Parsing " + str(len(files)) + " files...")
for filename in files:
    parsed = parseFile(filename)
    images_file.write(parsed.tobytes())
    label = int(filename[0])
    labels_file.write((label).to_bytes(1,byteorder="big"))
