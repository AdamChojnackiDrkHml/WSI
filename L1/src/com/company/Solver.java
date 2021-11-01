package com.company;
import java.util.ArrayList;
import java.util.Arrays;
import java.util.PriorityQueue;
import java.util.Stack;

import static com.company.Constants.BOARD_SIZE;
import static com.company.Constants.CURR_HEURISTIC;

public class Solver
{


    public Node getLastNode()
    {
        return lastNode;
    }

    private static class Node implements Comparable<Node>
    {
        private final byte[]  board_content;
        private byte heuristic;
        private final byte moves;
        private final Node prevNode;

        public Node(byte[]  board_content, byte moves, Node prev)
        {
            this.board_content = board_content;
            this.moves = moves;
            this.prevNode = prev;
            calculateHeuristic();
        }

        public int compareTo(Node that)
        {
            int thisPriority = this.moves + this.heuristic;
            int thatPriority = that.moves + that.heuristic;

            return Integer.compare(thisPriority, thatPriority);
        }

        private byte hammingHeuristic()
        {
            byte outOfPlace = 0 ;
            for(int i = 0 ; i < BOARD_SIZE ; i++ )
            {
                for(int j = 0 ; j < BOARD_SIZE ; j++ )
                {
                    if(board_content[Board.translate2Dto1D(i,j)] != 0)
                    {
                        if(board_content[Board.translate2Dto1D(i,j)]-1 != i*BOARD_SIZE+j)
                        {
                            outOfPlace++ ;
                        }
                    }
                }
            }

            return outOfPlace ;
        }

        private byte manhattanHeuristic()
        {
            byte heur = 0 ;
            for(int i = 0 ; i < BOARD_SIZE; i++ )
            {
                for(int j = 0; j < BOARD_SIZE; j++)
                {
                    if(board_content[Board.translate2Dto1D(i,j)] != 0)
                    {
                        heur += (byte)Math.abs(((board_content[Board.translate2Dto1D(i,j)]-1) / BOARD_SIZE) - (Board.translate2Dto1D(i,j) / BOARD_SIZE));
                        heur += (byte)Math.abs(((board_content[Board.translate2Dto1D(i,j)]-1) % BOARD_SIZE) - (Board.translate2Dto1D(i,j) % BOARD_SIZE));
                    }
                }

            }

            int prize = 2;
            int currSolved = 0;
            for(int i = 0; i < BOARD_SIZE; i++)
            {

                if(board_content[i] - 1 != i)
                {
                    break;
                }
                heur -= prize;
                currSolved++;
                prize += 2;
            }

            if(prize > 2 + BOARD_SIZE * 2)
            {
                for(int i = 0; i < BOARD_SIZE; i++)
                {

                    if(board_content[BOARD_SIZE + i] - 1 != BOARD_SIZE + i)
                    {
                        break;
                    }
                    heur -= prize;
                    currSolved++;
                    prize += 2;
                }

            }

            if(prize > 2 + BOARD_SIZE * 2 * 2 && BOARD_SIZE > 2)
            {
                for(int i = 0; i < BOARD_SIZE; i++)
                {

                    if(board_content[BOARD_SIZE * 2 + i] - 1 != BOARD_SIZE * 2 + i)
                    {
                        break;
                    }
                    heur -= prize;
                    currSolved++;
                    prize += 2;
                }

            }

            if(prize > 2 + BOARD_SIZE * 2 * 2 * 2 && BOARD_SIZE == 4)
            {
                for(int i = 0; i < BOARD_SIZE; i++)
                {

                    if(board_content[BOARD_SIZE * 3 + i] - 1 != BOARD_SIZE * 3 + i)
                    {
                        break;
                    }
                    heur -= prize;
                    currSolved++;
                    prize += 2;
                }

            }

            byte[] head = Board.translate1Dto2D(currSolved);
            byte zeroPos = 0;
            for(int i = 0; i < board_content.length; i++)
            {
                if(board_content[i] == 0)
                {
                    zeroPos = (byte)i;
                    break;
                }
            }
            byte[] zeroCords = Board.translate1Dto2D(zeroPos);
            heur += 20 * (byte)Math.abs(head[0] - zeroCords[0]);
            heur += 20 * (byte)Math.abs(head[1] - zeroCords[1]);
            return heur;
        }

        private void calculateHeuristic()
        {
            switch (CURR_HEURISTIC)
            {
                case MANHATTAN -> heuristic = manhattanHeuristic();
                case HAMMING -> heuristic = hammingHeuristic();
            }
        }
    }

    private Node lastNode;
    private boolean solvable;

    public Solver()
    {
        if(!isSolvable())
        {
            System.out.println("Nierozwiązywalne");
            return;
        }
        PriorityQueue<Node> pq = new PriorityQueue<>();
        pq.add(new Node(Board.getInstance().getBoardContent(), (byte) 0, null));
        while (true)
        {
            Node removed = pq.poll();
            if(removed == null)
            {
                break;
            }
            if (Board.isSolved(removed.board_content))
            {
                lastNode = removed;
                solvable = true;
                break;
            }


            ArrayList<byte[]> outcomes = Board.getOutcomes(removed.board_content);
            for (byte[] outcome : outcomes)
            {
                if (removed.prevNode != null && Arrays.equals(removed.prevNode.board_content, outcome))
                {
                    continue;
                }
                pq.add(new Node(outcome, (byte) (removed.moves + 1), removed));
            }
        }
        displaySolution();

    }

    public void displaySolution()
    {
        if(!solvable)
        {
            return;
        }
        Stack<byte[]> stack = new Stack<>();
        Node node = lastNode;

        while(node != null)
        {
            stack.push(node.board_content);
            node = node.prevNode;
        }

        int counter = 0;
        while(!stack.empty())
        {
            System.out.println(counter++);
            Board.printBoard(stack.pop());
        }
    }

    private int getInvCount()
    {
        int inv_count = 0;
        byte[] board_content = Board.getInstance().getBoardContent();
        for (int i = 0; i < BOARD_SIZE * BOARD_SIZE - 1; i++)
        {
            for (int j = i + 1; j < BOARD_SIZE * BOARD_SIZE; j++)
            {
                if (board_content[j] > 0 && board_content[i] > 0 && board_content[i] > board_content[j])
                    inv_count++;
            }
        }
        return inv_count;
    }

    public boolean isSolvable()
    {
        int invCount = getInvCount();

        // If grid is odd, return true if inversion
        // count is even.
        if (BOARD_SIZE % 2 == 1)
        {
            return (invCount % 2 == 0);
        }
        else     // grid is even
        {
            int xPos = 0;
            for(int i = 0; i < BOARD_SIZE * BOARD_SIZE; i++)
            {
                if(Board.getInstance().getBoardContent()[i] == 0)
                {
                    xPos = i;
                    break;
                }
            }

            if (xPos % 2 == 0)
                return (invCount % 2 == 0);
            else
                return (invCount % 2 == 1);
        }
    }
}