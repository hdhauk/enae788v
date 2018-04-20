import sys
import matplotlib.pyplot as plt
import matplotlib.path as mpath
from shapely.geometry.polygon import LinearRing, Polygon
import csv
import argparse
import re
import math


def rotate(origin, point, angle):
    """
    Rotate a point counterclockwise by a given angle around a given origin.

    The angle should be given in radians.
    """
    ox, oy = origin
    px, py = point

    qx = ox + math.cos(angle) * (px - ox) - math.sin(angle) * (py - oy)
    qy = oy + math.sin(angle) * (px - ox) + math.cos(angle) * (py - oy)
    return qx, qy


def draw_arrow(a, b):
    plt.arrow(a[0], a[1], b[0] - a[0], b[1] - a[1],
              head_width=1, length_includes_head=True, color='g')


parser = argparse.ArgumentParser(description='Plot RRT outputs.')
parser.add_argument('-o', '--obstacles', metavar='obstacles', nargs=1,
                    help='csv file of obstacles on form (x,y,r)',
                    default="H4_obstacles.txt")
parser.add_argument('-svg', help='save plot as svg',
                    default=False, action="store_true")
parser.add_argument('-png', help='save plot as png',
                    default=False, action="store_true")
args = parser.parse_args()

stdin = sys.stdin

# read start / goal coordinates
match = re.search(
    '^start=\[(\d+\.\d+),\s*(\d+\.\d+)\]\s*goal=\[(\d+\.\d+),\s*(\d+\.\d+),\s*(\d+\.\d+)\]',
    stdin.readline())
start = (float(match.group(1)), float(match.group(2)))
goal = (float(match.group(3)), float(match.group(4)), float(match.group(5)))


# Print any output before START_PATH
line = stdin.readline()
while line != "START_PATH\n":
    if line != "\n":
        print(line)
    line = stdin.readline()

edge_regex = """(-*\d+\.\d+)\s*,\s*(-*\d+\.\d+)\s*,\s*(-*\d+\.\d+)\s*,\s*(-*\d+\.\d+)"""
edge_regex2 = """(-*\d+\.\d+)\s*,\s*(-*\d+\.\d+)\s*,\s*(-*\d+\.\d+)\s*,\s*(-*\d+\.\d+)\s*,\s*(-*\d+\.\d+)"""
# Save path
line = stdin.readline()
path = []
while line != "END_PATH\n":
    match = re.search(edge_regex2, line)
    x1, y1, = float(match.group(1)), float(match.group(3))
    x2, y2, = float(match.group(2)), float(match.group(4))
    theta = float(match.group(5))
    path.append(((x1, y1), (x2, y2), theta))
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
    match = re.search(edge_regex, line)
    x1, y1, = float(match.group(1)), float(match.group(3))
    x2, y2, = float(match.group(2)), float(match.group(4))
    tree.append(((x1, y1), (x2, y2)))
    line = stdin.readline()

# Print any output after END_TREE
for l in stdin:
    if l != "":
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

# Start/goal
ax.add_artist(plt.Circle((start[0], start[1]), zorder=10,
                         radius=1.2, color='g', fill=True))
ax.add_artist(plt.Circle((goal[0], goal[1]), radius=goal[2], color='g'))


# Obstacles
for o in obstacles:
    ax.add_artist(plt.Circle((o[0], o[1]), radius=o[2]))

# Tree
for v in tree:
    plt.plot(v[0], v[1], color='k', marker='', markersize=1)

# Path
robot_path = [
    (0.0, 0.0),
    (0.1, 0.0),
    (0.2, 0.0),
    (0.3, 0.0),
    (0.4, 0.0),
    (0.5, 0.0),
    (0.6, 0.0),
    (0.7, 0.0),
    (0.8, 0.0),
    (0.9, 0.0),
    (1.0, 0.0),
    (-0.1, 0.0),
    (-0.2, 0.0),
    (-0.3, 0.0),
    (-0.4, 0.0),
    (-0.5, 0.0),
    (-0.6, 0.0),
    (-0.7, 0.0),
    (-0.8, 0.0),
    (-0.9, 0.0),
    (-1.0, 0.0),
    (0.2, 0.1),
    (0.2, 0.2),
    (0.1, 0.3),
    (0.05, 0.4),
    (0.0, 0.5),
    (0.2, -0.1),
    (0.2, -0.2),
    (0.1, -0.3),
    (0.05, -0.4),
    (0.0, -0.5),
    (-0.8, 0.1),
    (-0.8, 0.2),
    (-0.8, 0.3),
    (-0.8, -0.1),
    (-0.8, -0.2),
    (-0.8, -0.3)]

# for p in robot_path:
#     plt.scatter(p[0], p[1], marker='o', color='r')

for p in path:
    x = p[0]
    y = p[1]
    for r in robot_path:
        point = rotate((0, 0), r, p[2])
        plt.scatter(x[0]+point[0], y[0]+point[1],
                    marker='o', color='r', s=5, zorder=11)


plt.grid()
ax.set_aspect('equal')

if args.png:
    plt.savefig("start=({:.2f},{:.2f})_end=({:.2f},{:.2f}).png".format(
        start[0], start[1], goal[0], goal[1]), format="png", dpi=1000)
if args.svg:
    plt.savefig("start=({:.2f},{:.2f})_end=({:.2f},{:.2f}).svg".format(
        start[0], start[1], goal[0], goal[1]), format="svg")
if not (args.svg or args.png):
    plt.show()
