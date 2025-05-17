"""
This script generates large text files filled with random numbers for testing purposes.
It can be used to create test files of specified sizes in gigabytes, useful for
testing file processing performance and memory handling in the EMS system.

Usage:
    python very_big_file_generator.py

The script will generate a file named 'input.txt' with random numbers,
one per line. By default, it creates a 0.1GB (100MB) file.
"""

import random


def write_random(target_size_in_gb: float, filename: str = "test.txt") -> None:
    """
    Generate a text file filled with random numbers up to the specified size.

    Args:
        target_size_in_gb (float): Target file size in gigabytes
        filename (str, optional): Name of the output file. Defaults to "test.txt"

    The function writes random integers (1 to 1,000,000) to the file,
    one number per line, until the target file size is reached.
    Numbers are written in batches for better performance.
    """
    target_size = int(target_size_in_gb * 1024 * 1024 * 1024)
    size_written = 0
    buffer_size = 10_000  # Write 10,000 numbers per batch for efficiency

    with(open(filename, "w")) as file:
        while size_written < target_size:
            batch = "\n".join(str(random.randint(1, 1_000_000)) for _ in range(buffer_size))
            batch += "\n"
            file.write(batch)
            size_written += len(batch)


if __name__ == "__main__":
    write_random(0.1, "input.txt")  # Generates a 0.1GB (100MB) file named input.txt