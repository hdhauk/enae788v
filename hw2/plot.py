import sys
import matplotlib.pyplot as plt
import csv
import argparse
import re


def draw_arrow(a, b):
    plt.arrow(a[0], a[1], b[0] - a[0], b[1] - a[1],
              head_width=1, length_includes_head=True, color='g')


parser = argparse.ArgumentParser(description='Plot RRT outputs.')
parser.add_argument('-o', '--obstacles', metavar='obstacles', nargs=1,
                    help='csv file of obstacles on form (x,y,r)',
                    default="obstacles.txt")
args = parser.parse_args()

stdin = sys.stdin

# read start / goal coordinates
match = re.search(
    '^start=\[(\d+\.\d+),\s*(\d+\.\d+)\]\s*end=\[(\d+\.\d+),\s*(\d+\.\d+)\]',
    stdin.readline())
start = (float(match.group(1)), float(match.group(2)))
goal = (float(match.group(3)), float(match.group(4)))
print(start)
print(goal)

# Print any output before START_PATH
line = stdin.readline()
while line != "START_PATH\n":
    if line != "\n":
        print(line)
    line = stdin.readline()

# Save path
line = stdin.readline()
path = []
while line != "END_PATH\n":
    match = re.search(
        '(\d+\.\d+)\s*,\s*(\d+\.\d+)\s*,\s*(\d+\.\d+)\s*,\s*(\d+\.\d+)',
        line)
    x1, y1, = float(match.group(1)), float(match.group(3))
    x2, y2, = float(match.group(2)), float(match.group(4))
    path.append(((x1, y1), (x2, y2)))
    line = stdin.readline()

# Print any output before START_TREE
line = stdin.readline()
while line != "START_TREE\n":
    if line != "\n":
        print(line)
    line = stdin.readline()

# Save tree
tree = []
line = stdin.readline()
while line != "END_TREE\n":
    match = re.search(
        '(\d+\.\d+)\s*,\s*(\d+\.\d+)\s*,\s*(\d+\.\d+)\s*,\s*(\d+\.\d+)',
        line)
    x1, y1, = float(match.group(1)), float(match.group(3))
    x2, y2, = float(match.group(2)), float(match.group(4))
    tree.append(((x1, y1), (x2, y2)))
    line = stdin.readline()

# Print any output after END_TREE
for l in stdin:
    print(l)

# Read in obstacles
obstacle_file = open(args.obstacles)
csv_obstacles = csv.reader(obstacle_file)
obstacles = []
for row in csv_obstacles:
    x, y, r, = int(row[0]), int(row[1]), int(row[2])
    obstacles.append((x, y, r))
obstacle_file.close()


""" Plotting """
fig, ax = plt.subplots()
ax.set_xlim(0, 100)
ax.set_ylim(0, 100)
# Obstacles
for o in obstacles:
    ax.add_artist(plt.Circle((o[0], o[1]), radius=o[2]))

# Tree
for v in tree:
    plt.plot(v[0], v[1], color='k', marker='o')

# Path
for p in path:
    print(p[0], p[1])
    plt.plot(p[0], p[1], color='r', marker='o')

plt.show()
