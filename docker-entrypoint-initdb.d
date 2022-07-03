echo "Creating new user g..."
mysql -uroot -p"pass" -e "CREATE USER 'g'@'%' IDENTIFIED BY 'pass';"
echo "Granting privileges..."
mysql -uroot -p"pass" -e "GRANT ALL PRIVILEGES ON *.* TO 'g'@'%';"
mysql -uroot -p"pass" -e "FLUSH PRIVILEGES;"
echo "All done."