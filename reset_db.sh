#!/bin/bash

# 数据库重置脚本
# 用于解决唯一索引冲突问题

echo "正在重置数据库..."

# 提示用户输入数据库密码
echo "请输入MySQL root密码（或直接按回车跳过）："
read -s password

if [ -z "$password" ]; then
    # 无密码连接
    mysql -u root << EOF
DROP DATABASE IF EXISTS ai_course;
CREATE DATABASE ai_course CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
SHOW DATABASES LIKE 'ai_course';
EOF
else
    # 有密码连接
    mysql -u root -p"$password" << EOF
DROP DATABASE IF EXISTS ai_course;
CREATE DATABASE ai_course CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
SHOW DATABASES LIKE 'ai_course';
EOF
fi

if [ $? -eq 0 ]; then
    echo "数据库重置成功！"
    echo "现在可以重新启动应用程序。"
else
    echo "数据库重置失败，请检查MySQL连接。"
    echo ""
    echo "你也可以手动执行以下SQL命令："
    echo "DROP DATABASE IF EXISTS ai_course;"
    echo "CREATE DATABASE ai_course CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;"
fi