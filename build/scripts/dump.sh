#!/bin/bash
mysqldump -u root -p --host 127.0.0.1 --port 3306 --ssl-mode=REQUIRED library_dev > build/db/dump_file.sql
