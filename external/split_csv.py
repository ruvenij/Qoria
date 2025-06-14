import os
import csv

def split_csv(input_file, output_dir, lines_per_file=1000000):
    os.makedirs(output_dir, exist_ok=True)

    with open(input_file, 'r', newline='', encoding='utf-8') as infile:
        reader = csv.reader(infile)
        header = next(reader)

        file_count = 1
        rows = []
        for i, row in enumerate(reader, start=1):
            rows.append(row)
            if i % lines_per_file == 0:
                write_chunk(output_dir, file_count, rows)
                file_count += 1
                rows = []

        if rows:
            write_chunk(output_dir, file_count, rows)

def write_chunk(output_dir, index, rows):
    chunk_name = os.path.join(output_dir, f"chunk_{index}.csv")
    with open(chunk_name, 'w', newline='', encoding='utf-8') as outfile:
        writer = csv.writer(outfile)
        writer.writerows(rows)
    print(f"Wrote {chunk_name}")

if __name__ == "__main__":
    split_csv(
        input_file='data/GO_test_5m.csv',
        output_dir='data/split',
        lines_per_file=700_000
    )
