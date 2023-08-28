from pyspark.sql import SparkSession
from pyspark.sql.functions import *
from pyspark.sql.window import Window
import os
import time

os.environ['PYSPARK_SUBMIT_ARGS'] = '--packages org.apache.spark:spark-streaming-kafka-0-10_2.12:3.2.0,org.apache.spark:spark-sql-kafka-0-10_2.12:3.2.0 pyspark-shell'

spark = SparkSession \
    .builder \
    .appName("getDataFromPostgresAndAggregate") \
    .config("spark.jars", "./postgresql-42.6.0.jar") \
    .master("local[*]") \
    .getOrCreate()
    
spark.sparkContext.setLogLevel("ERROR")


while True:
    df = spark.read.format("jdbc").option("url", "jdbc:postgresql://localhost:5431/streaming") \
        .option("driver", "org.postgresql.Driver") \
        .option("dbtable", "user_activity") \
        .option("user", "streamuser") \
        .option("password", "streampass") \
        .load()
    
    # aggregations
    popular_sport = df.groupBy("location", "activity_type").count()
    w_sport = Window.partitionBy("location").orderBy(col("count").desc())
    popular_sport_agg = popular_sport.withColumn("row",row_number().over(w_sport)) \
        .filter(col("row") == 1).drop("row")
    popular_sport_agg = popular_sport_agg.withColumnRenamed("activity_type", "popular_sport")
    popular_sport_agg = popular_sport_agg.withColumnRenamed("count", "popular_sport_count")
    
    active_user = df.groupBy("location", "user_id").count()
    w_user = Window.partitionBy("location").orderBy(col("count").desc())
    active_user_agg = active_user.withColumn("row",row_number().over(w_user)) \
        .filter(col("row") == 1).drop("row")
    active_user_agg = active_user_agg.withColumnRenamed("location", "city")
    active_user_agg = active_user_agg.withColumnRenamed("user_id", "most_active_user")
    active_user_agg = active_user_agg.withColumnRenamed("count", "most_active_user_activity_count")

    joined_df = popular_sport_agg.join(
        active_user_agg, popular_sport_agg.location == active_user_agg.city, 'inner'
        )
    joined_df = joined_df.select("city", "popular_sport", "popular_sport_count", "most_active_user", "most_active_user_activity_count")

    #load back to postgres
    joined_df.write.format("jdbc").option("url", "jdbc:postgresql://localhost:5431/streaming") \
        .option("driver", "org.postgresql.Driver") \
        .option("dbtable", "city_stats") \
        .option("user", "streamuser") \
        .option("password", "streampass") \
        .mode("overwrite") \
        .save()
    
    print("Table overwritten.")

    time.sleep(10)

