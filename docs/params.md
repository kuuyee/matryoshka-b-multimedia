# 处理参数文档

> 上传/下载文件时，可以通过 query 对文件进行压缩、格式转换等控制

## 图片 (image)

| API          | Key    | Desc                                                         | Value              | Example                        |
| ------------ | ------ | ------------------------------------------------------------ | ------------------ | ------------------------------ |
| postFile     | format | 忽略自动探测和 content-type 中的格式，强制指定图片格式       | jpg,png,gif,webp   | `?format=webp`                 |
| retrieveFile | format | 返回指定的图片格式而不是默认格式                             | jpg,png,gif,webp   | `example.jpg?format=webp`      |
| postFile     | orig   | 维持原图格式，不进行转换和缩略图处理                         | 1,true,0,false     | `?orig=true`                   |
| postFile     | thumb  | 保存前对图片做缩略图处理，只有一个数字时作为dimY进行等比缩放 | \[Y\], \[X\]x\[Y\] | `?thumb=200`, `?thumb=600x800` |
| retrieveFile | size   | 对图片进行缩放，只有一个数字时作为dimY进行等比缩放           | \[Y\], \[X\]x\[Y\] | `?size=200`, `?size=600x800`   |