"""
This script validates the output of the external merge sort implementation.
It performs two main checks:
1. Verifies that the output file is correctly sorted in ascending order
2. Ensures that the input and output files have the same number of lines (no data loss)

Usage:
    python test_output.py

The script expects:
- input.txt: The original input file
- output.txt: The sorted output file from the merge sort
"""

def is_sorted(file_path: str) -> bool:
    """
    Check if a file contains numbers in ascending order.

    Args:
        file_path (str): Path to the file to check

    Returns:
        bool: True if file is sorted, False otherwise

    The function reads the file line by line to handle large files
    efficiently without loading the entire file into memory.
    Invalid lines (non-integer values) are skipped with a warning.
    """
    try:
        with open(file_path, "r") as file:
            prev = None
            for line in file:
                try:
                    current = int(line.strip())
                except ValueError:
                    print(f"Skipping invalid line: {line.strip()}")
                    continue

                if prev is not None and current < prev:
                    print(f"File is not sorted. Found {current} after {prev}.")
                    return False

                prev = current
        print("File is sorted.")
        return True
    except FileNotFoundError:
        print(f"File not found: {file_path}")
        return False
    except Exception as e:
        print(f"An error occurred: {e}")
        return False


def compare_line_count(file1: str, file2: str) -> str:
    """
    Compare the number of lines between two files.

    Args:
        file1 (str): Path to the first file (typically input file)
        file2 (str): Path to the second file (typically output file)

    Returns:
        str: A message describing the comparison result

    This function ensures that no data was lost during the sort process
    by verifying that the input and output files have the same number of lines.
    """
    result = ""
    try:
        with open(file1, 'r') as f1, open(file2, 'r') as f2:
            lines1 = sum(1 for _ in f1)  # Count lines in file1
            lines2 = sum(1 for _ in f2)  # Count lines in file2

        if lines1 > lines2:
            result = f"File '{file1}' has more lines ({lines1}) than file '{file2}' ({lines2})."
        elif lines1 < lines2:
            result = f"File '{file2}' has more lines ({lines2}) than file '{file1}' ({lines1})."
        else:
            result = f"Both files have the same number of lines ({lines1})."

    except FileNotFoundError as e:
        result = f"Error: {e}"
    except Exception as e:
        result = f"An unexpected error occurred: {e}"

    print(result)
    return result


if __name__ == "__main__":
    # File paths for validation
    input_path = "input.txt"
    output_path = "output.txt"

    # Perform validation checks
    print("Checking if output is sorted...")
    is_sorted(output_path)
    
    print("\nComparing input and output file sizes...")
    compare_line_count(input_path, output_path)