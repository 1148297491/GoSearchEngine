from flask import (Flask, render_template, request, flash, redirect, url_for, make_response, jsonify, send_from_directory)
from werkzeug.utils import secure_filename
import os
from datetime import timedelta
from predict import pred

ALLOWED_EXTENSIONS = set(["png","jpg","JPG","PNG", "bmp","jpeg"])

def is_allowed_file(filename):
    return '.' in filename and filename.rsplit('.', 1)[1] in ALLOWED_EXTENSIONS

app = Flask(__name__)
app.config['UPLOAD_FOLDER'] = '/home/GaoJH/disk0/lzm/pyt/findGraph/uploads/'#填绝对路径
#设置编码
app.config['JSON_AS_ASCII'] = False
# 静态文件缓存过期时间
app.send_file_max_age_default = timedelta(seconds=1)

@app.route('/uploads/<filename>')
def uploaded_file(filename):
	return send_from_directory(app.config['UPLOAD_FOLDER'],filename)

@app.route("/uploads",methods = ['POST', 'GET'])
def uploads():
	if request.method == "POST":
		#file = request.files['file']
		uploaded_files = request.files.getlist("file[]")
		urls = []
		for file in uploaded_files:
			if file and is_allowed_file(file.filename):
				filename = secure_filename(file.filename)
				file.save(os.path.join(app.config['UPLOAD_FOLDER'], filename))
				img_path = os.path.join(app.config['UPLOAD_FOLDER'], filename)
				urls = pred(img_path)
				# print(class_index,class_name)
				#data.append(result)
		#res_json = json.dumps({"status": "200", "msg": "success","data":data})
		return render_template('upload_more_ok.html',
								urls=urls,
								)
		
	return render_template("upload_more.html")


if __name__ == "__main__":
    # 0.0.0.0表示你监听本机的所有IP地址上,通过任何一个IP地址都可以访问到.
    # port为端口号
    # debug=Fasle表示不开启调试模式
    app.run(host='0.0.0.0', port=5000, debug=False)
