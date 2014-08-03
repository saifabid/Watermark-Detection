Watermark-Detection
===================

An app which will detect watermarks in images

Recently made this app in attempt to detect watermarks on images from an http API. Checks the top half, and bottom half of picture seperatly and writes the information for the watermarked images in JSON format to a file.

Uses both a cropping and OCR library
- github.com/oliamb/cutter was used for cropping the image
- github.com/otiai10/gosseract was used to run tesseract OCR on image 
