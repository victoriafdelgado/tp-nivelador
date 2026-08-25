from lottery.bet import Bet

def string_to_bet(bet_string):
    fields = bet_string.split(',')
    agency_id, first_name, last_name, document, birthdate, number  = fields
    return Bet(int(agency_id), first_name, last_name, int(document), birthdate, int(number))

def bet_to_string(bet):
    return f"{bet.first_name},{bet.last_name},{bet.document},{bet.birthdate},{bet.number}"

def strings_to_bets(payload):
    return [string_to_bet(line) for line in payload]