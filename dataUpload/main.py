import os
import json
import csv
from pymongo import MongoClient, InsertOne, errors
from dateutil import parser
from dotenv import load_dotenv

DATA_DIR = "../data/"
MAX_BATCH_SIZE = 50

# Colors
RED = "\033[91m"
YELLOW = "\033[93m"
GREEN = "\033[92m"
RESET = "\033[0m"

def log_message(prefix, message, color=RESET):
    print(f"{color}[{prefix}]{RESET} {message}")

def upload_documents(db, filepath, collection_name, row_parser):
    collection = db[collection_name]
    filepath = os.path.join(DATA_DIR, filepath)

    try:
        with open(filepath, newline='', encoding='utf-8') as csvfile:
            reader = csv.DictReader(csvfile)
            batch = []

            for idx, row in enumerate(reader):
                try:
                    doc = row_parser(row)
                    if doc is None:
                        log_message("WARN", f"Skipping invalid row at line {idx + 2}", YELLOW)
                        continue

                    batch.append(InsertOne(doc))

                    if len(batch) >= MAX_BATCH_SIZE:
                        collection.bulk_write(batch, ordered=False)
                        log_message("OK", f"Inserted {len(batch)} documents up to line {idx + 2}", GREEN)
                        batch.clear()

                except Exception as e:
                    log_message("ERROR", f"Unexpected error processing row {idx + 2}: {e}", RED)

            if batch:
                #collection.bulk_write(batch, ordered=False)
                log_message("OK", f"Inserted {len(batch)} remaining documents.", GREEN)

    except FileNotFoundError:
        log_message("ERROR", f"File {filepath} not found.", RED)
    except csv.Error as e:
        log_message("ERROR", f"CSV error while reading {filepath}: {e}", RED)
    except errors.BulkWriteError as bwe:
        log_message("ERROR", f"Bulk write error: {bwe.details}", RED)
    except Exception as e:
        log_message("ERROR", f"Unexpected error: {e}", RED)

def parse_user_row(row):
    try:
        return {
            "_id": int(row["_id"]),
            "Firstname": str(row["Firstname"]),
            "Lastname": str(row["Lastname"]),
            "Age": int(row["Age"]),
            "Gender": str(row["Gender"]),
        }
    except (KeyError, ValueError, TypeError) as e:
        print(f"{RED}[ERROR]{RESET} Parsing row {row}: {e}")
        return None

def parse_menu_items(menu_data):
    if 'Items' in menu_data:
        for item in menu_data['Items']:
            if '_id' in item:
                try:
                    item['_id'] = int(item['_id'])  # Convert _id to integer
                except ValueError as e:
                    print(f"{RED}[ERROR]{RESET} Invalid _id value: {e} for item: {item}")
    return menu_data

def parse_restaurant_row(row):
    try:
        coords_raw = row["Location"]
        menu_raw = json.loads(row["Menu"])

        coords = None
        if coords_raw:
            coords_list = json.loads(coords_raw)["Coordinates"]
            if isinstance(coords_list, list) and len(coords_list) == 2:
                coords = {
                    "type": "Point",
                    "coordinates": coords_list
                }

        menu = parse_menu_items(menu_raw)

        return {
            "_id": int(row["_id"]),
            "name": str(row["Name"]),
            "location": coords,
            "dob": parser.parse(row["Dob"]),
            "category": str(row["Category"]),
            "pricing": int(row["Pricing"]),
            "menu": menu,
        }

    except (KeyError, ValueError, TypeError) as e:
        print(f"{RED}[ERROR]{RESET} Parsing row {row}: {e}")
        return None

def parse_dishes_row(row):
    try:
        return {
            "_id": int(row["_id"]),
            "name": str(row["Name"]),
            "country": str(row["Country"]),
        }
    except (KeyError, ValueError, TypeError) as e:
        print(f"{RED}[ERROR]{RESET} Parsing row {row}: {e}")
        return None

def parse_reviews_row(row):
    try:
        return {
            "_id": int(row["_id"]),
            "restaurant": int(row["Restaurant"]),
            "client": int(row["Client"]),
            "rating": int(row["Rating"]),
            "relevance": int(row["Relevance"]),
        }
    except (KeyError, ValueError, TypeError) as e:
        print(f"{RED}[ERROR]{RESET} Parsing row {row}: {e}")
        return None
    

def parse_orders_row(row):
    try:
        item_raw = row["Item"]
        item = json.loads(item_raw)
        item["Price"] = float(item["Price"])

        return {
            "_id": int(row["_id"]),
            "client": int(row["Client"]),
            "restaurant": int(row["Restaurant"]),
            "state": str(row["State"]),
            "date": parser.parse(row["Date"]),
            "pricing": float(row["Pricing"]),
            "quantity": int(row["Quantity"]),
            "item": item,
        }

    except (KeyError, ValueError, TypeError) as e:
        print(f"{RED}[ERROR]{RESET} Parsing row {row}: {e}")
        return None


def main():
    load_dotenv()
    DB_URI = os.getenv("DB_URI")
    DB_NAME = os.getenv("DB_NAME")

    print(f"Connecting to DB {DB_URI}...")
    client = MongoClient(DB_URI)
    db = client[DB_NAME]

    print(f"============\nINSERTING USERS\n============")
    upload_documents(db, "../data/users.csv", "users", parse_user_row)

    print(f"============\nINSERTING RESTAURANTS\n============")
    upload_documents(db, "../data/restaurants.csv", "restaurants", parse_restaurant_row)

    print(f"============\nINSERTING DISHES\n============")
    upload_documents(db, "../data/dishes.csv", "dishes", parse_menu_items)

    print(f"============\nINSERTING REVIEWS\n============")
    upload_documents(db, "../data/reviews.csv", "reviews", parse_reviews_row)

    print(f"============\nINSERTING ORDERS\n============")
    upload_documents(db, "../data/orders.csv", "orders", parse_orders_row)

    print("✅ Done inserting all data!")

if __name__ == "__main__":
    main()