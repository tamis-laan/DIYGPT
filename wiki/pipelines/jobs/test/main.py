import logging
import sys

from pyflink.datastream import StreamExecutionEnvironment
from pyflink.table import StreamTableEnvironment


def main():
    # Create the execution environment
    env = StreamExecutionEnvironment.get_execution_environment()

    # Set number of execution threads to use
    env.set_parallelism(1)

    # Create a Stream Table environment
    t_env = StreamTableEnvironment.create(stream_execution_environment=env)

    # Create kafka source
    t_env.execute_sql("""
        CREATE TABLE source (
            meta ROW<
                uri      STRING,
                domain   STRING,
                topic STRING
            >,
            database STRING,
            page_id BIGINT,
            page_title STRING,
            page_namespace INT,
            rev_id BIGINT,
            rev_timestamp TIMESTAMP(3),
            rev_sha1 STRING,
            rev_minor_edit BOOLEAN,
            rev_len INT,
            rev_content_model STRING,
            rev_content_format STRING,
            performer ROW<
                user_text STRING,
                user_groups ARRAY<STRING>,
                user_is_bot BOOLEAN,
                user_id INT,
                user_edit_count BIGINT
            >,
            `page_is_redirect` BOOLEAN,
            `comment` STRING,
            `parsedcomment` STRING
        ) WITH ('connector' = 'kafka')
    """)

    # Create a print sink
    t_env.execute_sql("CREATE TABLE sink WITH ('connector' = 'print') LIKE source")

    # Set 
    t_env.execute_sql("""
        ALTER TABLE source SET (
          'topic' = 'wiki.page-create',
          'properties.bootstrap.servers' = 'kafka-cluster-kafka-bootstrap:9092',
          'properties.group.id' = 'flink.wiki.create-page',
          'scan.startup.mode' = 'earliest-offset',
          'json.fail-on-missing-field' = 'false',
          'json.ignore-parse-errors' = 'true',
          'format' = 'json'
    )""")

    # Move from source to sink 
    t_env.execute_sql("""
        INSERT INTO sink SELECT * FROM source""")


if __name__ == '__main__':
    logging.basicConfig(stream=sys.stdout, level=logging.INFO, format="%(message)s")
    main()
