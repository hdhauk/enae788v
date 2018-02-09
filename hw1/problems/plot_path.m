% Copyright 2018, Michael Otte
%
% Permission is hereby granted, free of charge, to any person obtaining a
% copy of this software and associated documentation files 
% (the "Software"), to deal in the Software without restriction, including
% without limitation the rights to use, copy, modify, merge, publish, 
% distribute, sublicense, and/or sell copies of the Software, and to permit
% persons to whom the Software is furnished to do so, subject to the 
% following conditions:
%
% The above copyright notice and this permission notice shall be included 
% in all copies or substantial portions of the Software.
%
%THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS 
% OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF 
% MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. 
% IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY 
% CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, 
% TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE 
% SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
%
% this will display the search tree and path
% assuming that the files have been generated
clc; close all; clear
p = '1';
search_tree_raw = csvread(strcat('search_tree',p,'.txt'));
path_raw = csvread(strcat('output_path',p,'.txt'));
nodes_raw = csvread(strcat('nodes_',p,'.txt'));
edges_raw = csvread(strcat('edges_' ,p,'.txt'));
start_goal = csvread(strcat('start_goal' ,p,'.txt'));

% a bit of data processing for faster plotting
search_tree = nan(3*size(search_tree_raw, 1), 2);

search_tree(1:3:end-2, 1) = search_tree_raw(:, 2);
search_tree(2:3:end-1, 1) = search_tree_raw(:, 5);
search_tree(1:3:end-2, 2) = search_tree_raw(:, 3);
search_tree(2:3:end-1, 2) = search_tree_raw(:, 6);

nodes = nodes_raw(2:end,2:3);

edges_raw = edges_raw(2:end,:);

edges = nan(3*size(edges_raw, 1), 2);

edges(1:3:end-2, 1) = nodes(edges_raw(:, 1),1);
edges(2:3:end-1, 1) = nodes(edges_raw(:, 2),1);
edges(1:3:end-2, 2) = nodes(edges_raw(:, 1),2);
edges(2:3:end-1, 2) = nodes(edges_raw(:, 2),2);


figure(1)
hold on
plot(nodes(:,1), nodes(:,2), 'ok')
plot(edges(:,1), edges(:,2), 'k')
plot(search_tree(:, 1), search_tree(:, 2), 'm', 'LineWidth', 2);
plot(path_raw(:,2), path_raw(:,3), 'b', 'LineWidth', 3);
plot(start_goal(1,1), start_goal(1,2),'or', 'LineWidth', 5);
plot(start_goal(2,1), start_goal(2,2),'og', 'LineWidth', 5);
hold off
print(strcat('solved_graph_',p,'.pdf'),'-dpdf');