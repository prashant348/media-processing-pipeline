import os
import shutil

def cleanup():

    target_dirs = { "output", "tmp" }
    deleted_count = 0

    cwd_path = os.getcwd()

    output = os.walk(cwd_path, topdown=False)
    for root, dirs, files in output:
        for name in dirs:
            if name in target_dirs:
                dir_path = os.path.join(root, name)
                print(f"Deleting {dir_path}")
                try:
                    shutil.rmtree(dir_path)
                    deleted_count += 1
                except Exception as e:
                    print(f"Error deleting {dir_path}: {e}")

    print(f"\nCleanup complete! Total {deleted_count} folders removed.")


if __name__ == "__main__":
    cleanup()