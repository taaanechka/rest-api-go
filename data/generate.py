import json, os.path

def mongo_init(init_dir, filename):
    data = [
        {
        "email": "dev.tati@mail.ru",
        "username": "tati",
        "password": "15827"
        },
        {
        "email": "dev.mr@mail.ru",
        "username": "mr",
        "password": "56397"
        }
    ]
    
    cur_dir = os.path.dirname(os.path.abspath(__file__))
    with open(f'{cur_dir}/{init_dir}{filename}', 'w', encoding='utf-8') as f:
        json.dump(data, f, ensure_ascii=False, indent=4)

if __name__ == "__main__":
    filename = 'data.json'
    init_dir_mongo = 'mongodb/init/'

    mongo_init(init_dir_mongo, filename)
