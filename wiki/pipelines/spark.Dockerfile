FROM spark:3.5.1-scala2.12-java11-python3-r-ubuntu

# Set root user
USER root

# Copy over requirements
COPY spark.requirements.txt requirements.txt

# Install requirements
RUN pip3 install -r requirements.txt

# Setup the app dir owned by spark user
RUN mkdir /app && \
	chown spark:spark /app && \
	chmod u=rwx,g=rwx,o=rx /app

# Set user spark
USER spark

# Set the working directory for the jobs
WORKDIR /app

# Add spark sql
RUN wget -P /opt/spark/jars/ https://repo1.maven.org/maven2/org/apache/spark/spark-sql_2.12/3.5.1/spark-sql_2.12-3.5.1.jar

# Get the kafka connector
RUN wget -P /opt/spark/jars/ https://repo1.maven.org/maven2/org/apache/spark/spark-sql-kafka-0-10_2.12/3.5.1/spark-sql-kafka-0-10_2.12-3.5.1.jar

RUN wget -P /opt/spark/jars/ https://repo1.maven.org/maven2/org/apache/spark/spark-token-provider-kafka-0-10_2.12/3.5.1/spark-token-provider-kafka-0-10_2.12-3.5.1.jar

# Get the kafka clients
RUN wget -P /opt/spark/jars/ https://repo1.maven.org/maven2/org/apache/kafka/kafka-clients/3.7.0/kafka-clients-3.7.0.jar

# Get commons pool
RUN wget -P /opt/spark/jars/ https://repo1.maven.org/maven2/org/apache/commons/commons-pool2/2.12.0/commons-pool2-2.12.0.jar

# Install postgresql client
RUN wget -P /opt/spark/jars/ https://repo1.maven.org/maven2/org/postgresql/postgresql/42.7.3/postgresql-42.7.3.jar

# Copy over the jobs
COPY jobs /app/jobs
