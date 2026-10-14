insert into players (id, name, primary_roles) values
  ('ratlab', 'หนูทดลอง', '["soft", "hard"]');

insert into ratings (subject_id, rater, traits, positions, scores) values (
  'ratlab',
  'เพื่อนร่วมทีม',
  '["summoner", "auraMan", "chameleon", "realSupport", "madman", "pest", "firstMover", "dogpiler", "cardio", "yesMan"]',
  '{
    "carry":   {"score": 1,  "comment": "เพราะไม่เคยเห็นเล่นเลย"},
    "mid":     {"score": 10, "comment": "เพราะฮีโร่มีจำกัดมาก ๆ เล่นได้แค่ตัว Summon แต่ถึงจะเล่นได้จำกัดแต่พอถูไถได้"},
    "offlane": {"score": 10, "comment": "เหตุผลเดียวกับ Midlane"},
    "hard":    {"score": 17, "comment": "เก่งจนน่ากลัว แต่ข้อจำกัดคือจำนวน Hero ที่เล่นได้"},
    "soft":    {"score": 17, "comment": "เหตุผลเดียวกับ Hard Support"}
  }',
  '{
    "skill":           {"score": 8,  "comment": "คะแนนที่หักไป เพราะว่าบางทีถ้าเล่นตัวถนัด แต่ไม่สามารถกดฝั่งตรงข้ามได้ อาจทำให้เกมเสียเพราะหนูทดลองฟาร์มน้อย ฟื้นไม่ได้ จมไปเลย"},
    "knowledge":       {"score": 8,  "comment": "ถูกหักเพราะเล่นฮีโร่เดิม ๆ ไม่มี Meta แต่เพราะเล่นฮีโร่ประเภทเดียวทำให้เก่งในความเข้าใจเกมที่มีต่อฮีโร่ตัวนั้น"},
    "consistency":     {"score": 10},
    "teamplay":        {"score": 8,  "comment": "ถ้าไม่ได้บอกไว้ หรือคุยทีมไว้ก่อน ถูกหัก 1 แต้ม เพราะจะรู้แค่ timing ของตัวเองเท่านั้น แต่จะไม่รู้ว่าเพื่อนต้องการ"},
    "laning":          {"score": 10},
    "lategame":        {"score": 8,  "comment": "หัก 1 คะแนนเพราะด้วยธรรมชาติของฮีโร่ทำให้เล่นได้ไม่สุดจนถึง Late game แล้วด้วยกับ Hero pool ที่เล่นได้จำกัด"},
    "captaincy":       {"score": 3,  "comment": "สั่งได้ แต่ไม่สามารถคุมทิศทางของเกมได้ 100%"},
    "communication":   {"score": 5},
    "roleFlexibility": {"score": 3,  "comment": "เพราะเล่น mid กับ hard/soft support ได้"},
    "heroPoolDepth":   {"score": 3,  "comment": "ถ้าทีมเดียวกัน ฮีโร่จำกัดเกินไป สำหรับ Captain Mode"}
  }'
);
