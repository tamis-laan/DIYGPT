FROM flink:1.18

# install python3: it has updated Python to 3.9 in Debian 11 and so install Python 3.7 from source, \
# it currently only supports Python 3.6, 3.7 and 3.8 in PyFlink officially.
RUN apt-get update -y && \
	apt-get install -y build-essential libssl-dev zlib1g-dev libbz2-dev libffi-dev && \
	wget https://www.python.org/ftp/python/3.7.9/Python-3.7.9.tgz && \
	tar -xvf Python-3.7.9.tgz && \
	cd Python-3.7.9 && \
	./configure --without-tests --enable-shared && \
	make -j6 && \
	make install && \
	ldconfig /usr/local/lib && \
	cd .. && rm -f Python-3.7.9.tgz && rm -rf Python-3.7.9 && \
	ln -s /usr/local/bin/python3 /usr/local/bin/python && \
	apt-get clean && \
	rm -rf /var/lib/apt/lists/*

# Switch to Flink user
USER flink

# Set the working directory for the jobs
WORKDIR /app

# Copy over requirements
COPY flink.requirements.txt requirements.txt

# Install requirements
RUN pip3 install -r requirements.txt

# Get the kafka connector
RUN wget -P /opt/flink/lib/ https://repo1.maven.org/maven2/org/apache/flink/flink-connector-kafka/3.1.0-1.18/flink-connector-kafka-3.1.0-1.18.jar

# Get the kafka clients
RUN wget -P /opt/flink/lib/ https://repo1.maven.org/maven2/org/apache/kafka/kafka-clients/3.7.0/kafka-clients-3.7.0.jar

# Install JDBC connector for sql databases
RUN wget -P /opt/flink/lib/ https://repo1.maven.org/maven2/org/apache/flink/flink-connector-jdbc/3.1.2-1.18/flink-connector-jdbc-3.1.2-1.18.jar

# Install postgresql client
RUN wget -P /opt/flink/lib/ https://repo1.maven.org/maven2/org/postgresql/postgresql/42.7.3/postgresql-42.7.3.jar

# Install connector for files
RUN wget -P /opt/flink/lib/ https://repo1.maven.org/maven2/org/apache/flink/flink-connector-files/1.18.1/flink-connector-files-1.18.1.jar

# Copy over the jobs
COPY jobs /app/jobs
