
import sys

class Node:
    def __init__(self, symbol, leftChild=None, rightChild=None):
        self.symbol = symbol
        self.leftChild = leftChild
        self.rightChild = rightChild


class Parser:
    def __init__(self):
        self.error = False
        self.next_token = '%'
        self.file = None

    def parse_file(self, file_path):
        try:
            self.file = open(file_path, "r")
            self.lex()
            self.parse()
            self.file.close()
        except IOError:
            print(f"Error: Failed to open file '{file_path}'.")

    def lex(self):
        while True:
            next_char = self.file.read(1)
            if next_char == ' ' or next_char == '\n':
                continue
            self.next_token = next_char
            break

    def unconsumed_input(self):
        return self.file.read()

    def parse(self):
        tree = self.G()
        if self.next_token == '$' and not self.error:
            print("success")
            self.print_tree(tree)
            value = self.evaluate(tree)
            print("The value is", value)
        else:
            print("Failure: unconsumed input =", self.unconsumed_input())

    def G(self):
        tree = self.E()
        if self.next_token == '$' and not self.error:
            return tree
        else:
            self.error = True
            return None

    def E(self):
        if self.error:
            return None
        print("E -> T R")
        temp = self.T()
        return self.R(temp)

    def R(self, tree):
        if self.error:
            return None
        if self.next_token == '+':
            print("R -> + T R")
            self.lex()
            temp1 = self.T()
            temp2 = self.R(temp1)
            return Node('+', tree, temp2)
        elif self.next_token == '-':
            print("R -> - T R")
            self.lex()
            temp1 = self.T()
            temp2 = self.R(temp1)
            return Node('-', tree, temp2)
        else:
            print("R -> e")
            return tree

    def T(self):
        if self.error:
            return None
        print("T -> F S")
        temp = self.F()
        return self.S(temp)

    def S(self, tree):
        if self.error:
            return None
        if self.next_token == '*':
            print("S -> * F S")
            self.lex()
            temp1 = self.F()
            temp2 = self.S(temp1)
            return Node('*', tree, temp2)
        elif self.next_token == '/':
            print("S -> / F S")
            self.lex()
            temp1 = self.F()
            temp2 = self.S(temp1)
            return Node('/', tree, temp2)
        else:
            print("S -> e")
            return tree

    def F(self):
        if self.error:
            return None
        if self.next_token == '(':
            print("F -> ( E )")
            self.lex()
            temp = self.E()
            if self.next_token == ')':
                self.lex()
                return temp
            else:
                self.error = True
                print("Error: Unexpected token", self.next_token)
                print("Unconsumed input", self.unconsumed_input())
                return None
        elif self.next_token in ['a', 'b', 'c', 'd']:
            print("F -> M")
            return self.M()
        elif self.next_token in ['0', '1', '2', '3']:
            print("F -> N")
            return self.N()
        else:
            self.error = True
            print("Error: Unexpected token", self.next_token)
            print("Unconsumed input", self.unconsumed_input())
            return None
        
    def N(self):
        prev_token = self.next_token
        if self.error:
            return None
        if self.next_token in ['0', '1', '2', '3']:
            print("N ->", self.next_token)
            self.lex()
            return Node(prev_token)
        else:
            self.error = True
            print("Error: Unexpected token", self.next_token)
            print("Unconsumed input", self.unconsumed_input())
            return None

    def M(self):
        prev_token = self.next_token
        if self.error:
            return None
        if self.next_token in ['a', 'b', 'c', 'd']:
            print("M ->", self.next_token)
            self.lex()
            return Node(prev_token)
        else:
            self.error = True
            return None

    def print_tree(self, tree, level=0):
        if tree is not None:
            print(' ' * level, tree.symbol)
            self.print_tree(tree.leftChild, level + 1)
            self.print_tree(tree.rightChild, level + 1)

    def evaluate(self, tree):
        if tree is None:
            return None
        if tree.symbol.isalpha():
            return self.get_variable_value(tree.symbol)
        if tree.symbol.isnumeric():
            return int(tree.symbol)
        left_value = self.evaluate(tree.leftChild)
        right_value = self.evaluate(tree.rightChild)
        if tree.symbol == '+':
            return left_value + right_value
        elif tree.symbol == '-':
            return left_value - right_value
        elif tree.symbol == '*':
            return left_value * right_value
        elif tree.symbol == '/':
            if right_value == 0:
                print("Error: Division by zero.")
                return None
            return left_value / right_value
        return None

    def get_variable_value(self, variable):
        # Implement your logic to retrieve variable values here
        # For example, you can use a dictionary or a database
        variables = {'a': 10, 'b': 5, 'c': 7, 'd': 2}
        return variables.get(variable, None)


if len(sys.argv) < 2:
        print("Usage: python3 run.py <file_path>")
        sys.exit(1)

file_path = sys.argv[1]

parser = Parser()
parser.parse_file(file_path)