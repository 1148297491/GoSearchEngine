# -*- coding: utf-8 -*-

import torch

from torchvision import transforms
from torchvision import models

import numpy as np 
import faiss
from PIL import Image

from PIL import ImageFile

class_names = ['bird', 'book', 'butterfly', 'cattle', 'chicken', 'elephant', 'horse', 
               'phone', 'sheep', 'shoes', 'spider', 'squirrel', 'watch']
ImageFile.LOAD_TRUNCATED_IMAGES = True
Image.MAX_IMAGE_PIXELS = None
transform = transforms.Compose([
    transforms.Resize((256, 256)),
    transforms.CenterCrop(224),
    transforms.ToTensor(),
    transforms.Normalize([0.485, 0.456, 0.406], [0.229, 0.224, 0.225])
])
device = torch.device('cuda:0' if torch.cuda.is_available() else 'cpu')
# 加载预训练模型
def load_model():
    model = models.resnet34(pretrained=True)
    model.to(device)
    model.eval()
    return model
# 定义 特征提取器
def feature_extract(model, x):
    x = model.conv1(x)
    x = model.bn1(x)
    x = model.relu(x)
    x = model.maxpool(x)
    x = model.layer1(x)
    x = model.layer2(x)
    x = model.layer3(x)
    x = model.layer4(x)
    x = model.avgpool(x)
    x = torch.flatten(x, 1)
    return x
fake_headers = {
    'User-Agent': 'Mozilla/5.0 (Windows NT 10.0; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/79.0.3945.79 Safari/537.36'
}
model = load_model()
index = faiss.read_index('index_file.index')
topK = 20
with open('wukong_urls.txt', 'r', encoding='utf-8') as f:
    urlsTotal = [row.strip() for row in f.readlines()]

def index_search(feat,topK ):
    feat = np.expand_dims( np.array(feat),axis=0 )
    feat = feat.astype('float32')
    
    dis,ind = index.search( feat,topK )
    return dis,ind # 距离，相似图片id

def visual_plot(ind):       

    idx = 0
    urls = []
    for idx in range(20):
        _id = ind[0][idx]
        urls.append(urlsTotal[_id])
    return urls

def pred(img_path):
    img = Image.open(img_path)
    img = transform(img) # torch.Size([3, 224, 224])
    img = img.unsqueeze(0) # torch.Size([1, 3, 224, 224])
    img = img.to(device)
    # 对我们的图片进行预测
    with torch.no_grad():
        # 图片-> 图片特征
        feature = feature_extract( model,img )
        # 特征-> 检索
        feature_list = feature.data.cpu().tolist()[0]
        # 相似图片可视化
        dis,ind = index_search( feature_list,topK=topK )
        return visual_plot( ind)
