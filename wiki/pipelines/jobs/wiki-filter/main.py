from pyspark.sql import SparkSession
from pyspark.sql.types import StructType, StructField, StringType, BooleanType, LongType, ArrayType, TimestampType
from pyspark.sql.functions import from_json, to_json, struct, col

schema = StructType([
    StructField("$schema", StringType(), True),
    StructField("meta", StructType([
        StructField("uri", StringType(), True),
        StructField("request_id", StringType(), True),
        StructField("id", StringType(), True),
        StructField("dt", TimestampType(), True),
        StructField("domain", StringType(), True),
        StructField("stream", StringType(), True),
        StructField("topic", StringType(), True),
        StructField("partition", LongType(), True),
        StructField("offset", LongType(), True)
    ]), True),
    StructField("database", StringType(), True),
    StructField("page_id", LongType(), True),
    StructField("page_title", StringType(), True),
    StructField("page_namespace", LongType(), True),
    StructField("rev_id", LongType(), True),
    StructField("rev_timestamp", TimestampType(), True),
    StructField("rev_sha1", StringType(), True),
    StructField("rev_minor_edit", BooleanType(), True),
    StructField("rev_len", LongType(), True),
    StructField("rev_content_model", StringType(), True),
    StructField("rev_content_format", StringType(), True),
    StructField("performer", StructType([
        StructField("user_text", StringType(), True),
        StructField("user_groups", ArrayType(StringType(), True), True),
        StructField("user_is_bot", BooleanType(), True),
        StructField("user_id", LongType(), True),
        StructField("user_registration_dt", TimestampType(), True),
        StructField("user_edit_count", LongType(), True)
    ]), True),
    StructField("page_is_redirect", BooleanType(), True),
    StructField("comment", StringType(), True),
    StructField("parsedcomment", StringType(), True),
    StructField("dt", TimestampType(), True),
    StructField("rev_slots", StructType([
        StructField("main", StructType([
            StructField("rev_slot_content_model", StringType(), True),
            StructField("rev_slot_sha1", StringType(), True),
            StructField("rev_slot_size", LongType(), True),
            StructField("rev_slot_origin_rev_id", LongType(), True)
        ]), True)
    ]), True)
])

if __name__ == "__main__":

    spark = SparkSession\
        .builder\
        .appName("wiki-filter")\
        .getOrCreate()

    # Set log level
    spark.sparkContext.setLogLevel("WARN")

    # Read from spark
    source = spark \
      .readStream \
      .format("kafka") \
      .option("kafka.bootstrap.servers", "kafka-cluster-kafka-bootstrap:9092") \
      .option("subscribe", "wiki.page-create") \
      .option("startingOffsets", "earliest") \
      .load()

    # Unpack json using schema
    tableDF = source.selectExpr("CAST(value AS STRING) as json_str") \
        .select(from_json("json_str", schema).alias("data")) \
        .select("data.*")

    # Print the schema
    tableDF.printSchema()

    # Filter out relevant data
    filteredDF = tableDF.select( 
        col("$schema").alias("schema"),
        col("meta.uri"),
        col("meta.domain"),
        col("meta.stream"),
        col("database"),
        col("page_id"),
        col("page_title"),
        col("page_namespace"),
        col("rev_content_model"),
        col("rev_content_format"),
        col("performer.user_is_bot"),
        col("performer.user_edit_count"),
        col("comment"),
        col("parsedcomment"),
        col("dt"),
    )

    # Write to console
    printer = filteredDF \
        .writeStream \
        .outputMode("append") \
        .format("console") \
        .trigger(processingTime='5 seconds') \
        .start()

    # Encode into json for kafka
    kafkaDF = filteredDF.select(to_json(struct([filteredDF[c] for c in filteredDF.columns])).alias("value"))

    # Write results to kafka
    sink = kafkaDF \
        .writeStream \
        .format('kafka') \
        .option("kafka.bootstrap.servers", "kafka-cluster-kafka-bootstrap:9092") \
        .option("topic", "wiki.page-create.filtered") \
        .option("checkpointLocation", "/app/checkpoint/test/") \
        .start()

    # Keep driver alive
    sink.awaitTermination()
    printer.awaitTermination()

