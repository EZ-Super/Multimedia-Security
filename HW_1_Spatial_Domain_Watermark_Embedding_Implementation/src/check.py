
from PIL import Image
import numpy as np

# 讀取圖片

image = Image.open("C:\\Users\\User\\Desktop\\github\\Multimedia Security\\Multimedia-Security\\HW_1_Spatial_Domain_Watermark_Embedding_Implementation\\images\\elaine_512x512.bmp").convert("RGB")
ImageArray = np.array(image)
# 檢查 mode
mode = image.mode
bit_depth = {
    "1": 1,     # 1-bit 黑白
    "L": 8,     # 8-bit 灰階
    "P": 8,     # 8-bit 調色盤
    "RGB": 24,  # 8-bit * 3 (每個通道 8-bit)
    "RGBA": 32, # 8-bit * 4 (每個通道 8-bit + Alpha)
    "LA": 16,   # 8-bit * 2 (灰階 + Alpha)
    "CMYK": 32, # 8-bit * 4 (印刷模式)
    "I": 32,    # 32-bit 整數
    "F": 32,    # 32-bit 浮點數
}.get(mode, "Unknown")

print(f"圖片模式: {mode}")
print(f"每個像素的位元數: {bit_depth} bit")

r,g,b = ImageArray[:,:,0], ImageArray[:,:,1], ImageArray[:,:,2]

if np.all(r == g) and np.all(g == b):
    # 找出所有唯一的顏色值
    unique_colors = np.unique(ImageArray.reshape(-1, 3), axis=0)

    # 判斷是否為黑白圖片
    if len(unique_colors) == 2 and \
            all((color == [0, 0, 0]).all() or (color == [255, 255, 255]).all() for color in unique_colors):
        print("這是一張 **黑白圖片（Binary Image）**")
    else:
        print("這是一張 **灰階圖片（Grayscale）**")
else:
    print("這是一張 **彩色圖片（Colored Image）**")


for x in range(image.width):
    for y in range(image.height):
        r,g,b = image.getpixel((x, y))
        print(f"{r} {g} {b}")