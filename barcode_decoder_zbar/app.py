import cv2
from pyzbar.pyzbar import decode
   
def BarcodeReader(image):
    print(f"read file ='{image}'")
    # read the image in numpy array using cv2
    img = cv2.imread(image)
       
    # Decode the barcode image
    detectedBarcodes = decode(img)
       
    # If barcode not detected
    if not detectedBarcodes:
        raise Exception("Barcode Not Detected or your barcode is blank/corrupted!")

    # Traveres through all the detected barcodes in image
    for barcode in detectedBarcodes:  
        
        # Locate the barcode position in image
        (x, y, w, h) = barcode.rect
        
        # print(x, y, w, h)
        # Put the rectangle in image using cv2 to heighlight the barcode
        # cv2.rectangle(img, (x-10, y-10),
        #               (x + w+10, y + h+10), 
        #               (255, 0, 0), 2)
            
        if barcode.data !="":
            print(f"value={barcode.data}. type={barcode.type}")
                  
  
if __name__ == "__main__":
    # image="barcodes.png"
    image="barcodes-no-text.png"
    BarcodeReader(image)